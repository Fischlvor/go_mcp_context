#!/usr/bin/env python3
"""
MCP 项目文档更新性能分析脚本
分析 activity_logs 表中的历史导入记录，生成详细的性能报告
"""

import subprocess
import json
import sys
from datetime import datetime
from collections import defaultdict
from db_config import get_sql_command

def run_sql(sql):
    """执行 SQL 查询并返回结果"""
    cmd = get_sql_command(sql)
    result = subprocess.run(cmd, shell=True, capture_output=True, text=True)
    if result.returncode != 0:
        print(f"SQL 执行失败: {result.stderr}", file=sys.stderr)
        return []
    return [line.strip() for line in result.stdout.strip().split('\n') if line.strip()]

def parse_timestamp(ts_str):
    """解析时间戳"""
    try:
        # 处理带时区的时间戳
        return datetime.strptime(ts_str.split('+')[0], '%Y-%m-%d %H:%M:%S.%f')
    except:
        try:
            return datetime.strptime(ts_str.split('+')[0], '%Y-%m-%d %H:%M:%S')
        except:
            return None

def get_all_tasks():
    """获取所有任务的基本信息"""
    sql = """
    SELECT 
        task_id, 
        COUNT(*) as event_count, 
        MIN(created_at) as start_time, 
        MAX(created_at) as end_time,
        EXTRACT(EPOCH FROM (MAX(created_at) - MIN(created_at))) as duration_seconds
    FROM activity_logs 
    WHERE task_id IS NOT NULL AND task_id != '' 
    GROUP BY task_id 
    ORDER BY start_time DESC 
    """
    lines = run_sql(sql)
    tasks = []
    for line in lines:
        parts = line.split('|')
        if len(parts) >= 5:
            tasks.append({
                'task_id': parts[0],
                'event_count': int(parts[1]),
                'start_time': parts[2],
                'end_time': parts[3],
                'duration_seconds': float(parts[4]) if parts[4] else 0
            })
    return tasks

def get_task_details(task_id):
    """获取任务的详细事件"""
    sql = f"""
    SELECT 
        al.library_id,
        l.name as library_name,
        al.event,
        al.status,
        al.message,
        al.created_at,
        al.version
    FROM activity_logs al 
    LEFT JOIN libraries l ON al.library_id = l.id 
    WHERE al.task_id = '{task_id}' 
    ORDER BY al.created_at
    """
    lines = run_sql(sql)
    events = []
    for line in lines:
        parts = line.split('|')
        if len(parts) >= 7:
            events.append({
                'library_id': parts[0],
                'library_name': parts[1],
                'event': parts[2],
                'status': parts[3],
                'message': parts[4],
                'created_at': parts[5],
                'version': parts[6] if parts[6] else 'latest'
            })
    return events

def analyze_task(task_id, events):
    """分析单个任务的性能"""
    if not events:
        return None
    
    analysis = {
        'task_id': task_id,
        'library_name': events[0]['library_name'],
        'version': events[0]['version'],
        'total_events': len(events),
        'start_time': events[0]['created_at'],
        'end_time': events[-1]['created_at'],
        'phases': defaultdict(list),
        'documents': defaultdict(dict),
        'operation_type': 'unknown'
    }
    
    # 确定操作类型
    if any('github.import' in e['event'] for e in events):
        analysis['operation_type'] = 'GitHub Import'
    elif any('version.refresh' in e['event'] for e in events):
        analysis['operation_type'] = 'Version Refresh'
    elif any('document.upload' in e['event'] for e in events):
        analysis['operation_type'] = 'Document Upload'
    
    # 分析每个阶段
    prev_time = None
    prev_event_type = None
    current_doc = None
    
    for event in events:
        event_type = event['event']
        timestamp = parse_timestamp(event['created_at'])
        
        if not timestamp:
            continue
        
        if prev_time:
            duration = (timestamp - prev_time).total_seconds()
        else:
            duration = 0
        
        # 提取文档名
        if 'document.' in event_type:
            msg = event['message']
            if ':' in msg:
                doc_name = msg.split(':')[1].strip().split('(')[0].strip()
                current_doc = doc_name
        
        # 修正阶段归属：document.complete 的耗时实际上是 Embedding 生成的时间
        # 因为 document.embed 事件是瞬时记录的，真正的 API 调用在之后
        actual_event_type = event_type
        if event_type == 'document.complete' and prev_event_type == 'document.embed':
            # 将这段时间归类为 Embedding 生成
            actual_event_type = 'document.embed.actual'
        
        # 记录阶段
        phase_info = {
            'event': actual_event_type,
            'message': event['message'],
            'timestamp': event['created_at'],
            'duration': duration,
            'document': current_doc
        }
        
        analysis['phases'][actual_event_type].append(phase_info)
        prev_event_type = event_type
        
        # 按文档分组
        if current_doc:
            if current_doc not in analysis['documents']:
                analysis['documents'][current_doc] = {
                    'events': [],
                    'total_duration': 0,
                    'chunks': 0,
                    'tokens': 0
                }
            analysis['documents'][current_doc]['events'].append(phase_info)
            analysis['documents'][current_doc]['total_duration'] += duration
            
            # 提取块数和 token 数
            if 'chunks' in event['message'] or '块' in event['message']:
                try:
                    chunks = int(event['message'].split('(')[1].split('块')[0].strip())
                    analysis['documents'][current_doc]['chunks'] = chunks
                except:
                    pass
            if 'tokens' in event['message']:
                try:
                    tokens = int(event['message'].split('tokens')[0].split(',')[-1].strip())
                    analysis['documents'][current_doc]['tokens'] = tokens
                except:
                    pass
        
        prev_time = timestamp
    
    # 计算总时长
    start = parse_timestamp(analysis['start_time'])
    end = parse_timestamp(analysis['end_time'])
    if start and end:
        analysis['total_duration'] = (end - start).total_seconds()
    else:
        analysis['total_duration'] = 0
    
    return analysis

def print_report(analyses):
    """打印性能报告"""
    print("=" * 100)
    print("MCP 项目文档更新性能分析报告")
    print("=" * 100)
    print()
    
    # 按操作类型分组
    by_type = defaultdict(list)
    for analysis in analyses:
        if analysis:
            by_type[analysis['operation_type']].append(analysis)
    
    for op_type, tasks in by_type.items():
        print(f"\n{'=' * 100}")
        print(f"操作类型: {op_type} (共 {len(tasks)} 个任务)")
        print(f"{'=' * 100}\n")
        
        for i, analysis in enumerate(tasks, 1):  # 只显示前 10 个
            print(f"\n{'─' * 100}")
            print(f"任务 #{i}: {analysis['library_name']} @ {analysis['version']}")
            print(f"{'─' * 100}")
            print(f"Task ID:      {analysis['task_id']}")
            print(f"开始时间:     {analysis['start_time']}")
            print(f"结束时间:     {analysis['end_time']}")
            print(f"总耗时:       {analysis['total_duration']:.2f} 秒")
            print(f"事件总数:     {analysis['total_events']}")
            # 统计 document.parse 事件数作为文档数量
            doc_parse_count = len(analysis['phases'].get('document.parse', []))
            print(f"文档数量:     {doc_parse_count} 个 (document.parse 事件数)")
            
            # 阶段统计
            print(f"\n  阶段统计:")
            phase_stats = {}
            for phase, events in analysis['phases'].items():
                total_duration = sum(e['duration'] for e in events)
                phase_stats[phase] = {
                    'count': len(events),
                    'total_duration': total_duration,
                    'avg_duration': total_duration / len(events) if events else 0
                }
            
            # 合并 document.embed 和 document.embed.actual 的统计
            if 'document.embed.actual' in phase_stats:
                if 'document.embed' not in phase_stats:
                    phase_stats['document.embed'] = {'count': 0, 'total_duration': 0, 'avg_duration': 0}
                phase_stats['document.embed']['count'] += phase_stats['document.embed.actual']['count']
                phase_stats['document.embed']['total_duration'] += phase_stats['document.embed.actual']['total_duration']
                if phase_stats['document.embed']['count'] > 0:
                    phase_stats['document.embed']['avg_duration'] = phase_stats['document.embed']['total_duration'] / phase_stats['document.embed']['count']
                del phase_stats['document.embed.actual']
            
            # 按耗时排序
            sorted_phases = sorted(phase_stats.items(), key=lambda x: x[1]['total_duration'], reverse=True)
            for phase, stats in sorted_phases[:10]:  # 只显示前 10 个阶段
                print(f"    • {phase:30s} - 次数: {stats['count']:3d}, 总耗时: {stats['total_duration']:7.2f}s, 平均: {stats['avg_duration']:6.2f}s")
    
    # 总体统计
    print(f"\n\n{'=' * 100}")
    print("总体统计")
    print(f"{'=' * 100}\n")
    
    all_durations = [a['total_duration'] for a in analyses if a and a['total_duration'] > 0]
    # 使用 document.parse 事件数统计文档数量
    all_doc_counts = [len(a['phases'].get('document.parse', [])) for a in analyses if a]
    
    if all_durations:
        print(f"任务总数:         {len(analyses)}")
        print(f"平均耗时:         {sum(all_durations) / len(all_durations):.2f} 秒")
        print(f"最短耗时:         {min(all_durations):.2f} 秒")
        print(f"最长耗时:         {max(all_durations):.2f} 秒")
        print(f"平均文档数:       {sum(all_doc_counts) / len(all_doc_counts):.1f}")
        print()

def main():
    print("正在获取任务列表...")
    tasks = get_all_tasks()
    print(f"找到 {len(tasks)} 个任务\n")
    
    print("正在分析任务详情...")
    analyses = []
    for i, task in enumerate(tasks, 1):
        print(f"\r  [{i}/{len(tasks)}] 分析任务 {task['task_id']}...", end='', flush=True)
        events = get_task_details(task['task_id'])
        # 检查是否有成功完成的事件（status='success'）
        has_success_event = any(e.get('status') == 'success' for e in events)
        
        analysis = analyze_task(task['task_id'], events)
        # 过滤掉无效任务：必须有 success 状态且耗时>0
        if analysis and has_success_event and analysis['total_duration'] > 0:
            analyses.append(analysis)
    
    print("\n\n正在生成报告...\n")
    print_report(analyses)

if __name__ == '__main__':
    main()
