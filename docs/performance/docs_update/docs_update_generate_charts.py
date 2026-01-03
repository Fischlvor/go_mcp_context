#!/usr/bin/env python3
"""
生成性能分析图表
"""

import subprocess
import matplotlib.pyplot as plt
import seaborn as sns
import numpy as np
from collections import defaultdict
import json
from db_config import get_sql_command

# 设置中文字体
import matplotlib
import matplotlib.font_manager as fm
from matplotlib.font_manager import FontProperties

# 查找系统中的中文字体文件（优先 ttf，避免 ttc）
result = subprocess.run(['fc-list', ':lang=zh', 'file'], capture_output=True, text=True)
chinese_font_files = []
for line in result.stdout.strip().split('\n'):
    if line and ('.ttf' in line.lower() or '.otf' in line.lower()):
        font_file = line.split(':')[0].strip()
        # 优先选择 Noto 或 AR PL 字体
        if 'Noto' in font_file or 'AR' in font_file or 'ukai' in font_file or 'uming' in font_file:
            chinese_font_files.insert(0, font_file)
        else:
            chinese_font_files.append(font_file)

# 全局字体属性
CHINESE_FONT = None
if chinese_font_files:
    font_path = chinese_font_files[0]
    print(f"✓ 使用中文字体文件: {font_path}")
    CHINESE_FONT = FontProperties(fname=font_path)
    # 设置全局默认字体
    matplotlib.rcParams['font.family'] = CHINESE_FONT.get_name()
else:
    print("⚠ 未找到中文字体，图表中文可能显示为方框")

matplotlib.rcParams['axes.unicode_minus'] = False

# 设置 seaborn 样式
sns.set_style("whitegrid")
sns.set_palette("husl")

def run_sql(sql):
    """执行 SQL 查询并返回结果"""
    cmd = get_sql_command(sql)
    result = subprocess.run(cmd, shell=True, capture_output=True, text=True)
    if result.returncode != 0:
        print(f"SQL 执行失败: {result.stderr}")
        return []
    return [line.strip() for line in result.stdout.strip().split('\n') if line.strip()]

def get_task_summary():
    """获取任务汇总数据（只包含成功的任务）"""
    sql = """
    SELECT 
        task_id, 
        COUNT(*) as event_count, 
        MIN(created_at) as start_time, 
        MAX(created_at) as end_time,
        EXTRACT(EPOCH FROM (MAX(created_at) - MIN(created_at))) as duration_seconds,
        MAX(CASE WHEN status = 'success' THEN 1 ELSE 0 END) as has_success
    FROM activity_logs 
    WHERE task_id IS NOT NULL AND task_id != '' 
    GROUP BY task_id 
    HAVING MAX(CASE WHEN status = 'success' THEN 1 ELSE 0 END) = 1
       AND EXTRACT(EPOCH FROM (MAX(created_at) - MIN(created_at))) > 0
    ORDER BY start_time DESC 
    LIMIT 50
    """
    lines = run_sql(sql)
    tasks = []
    for line in lines:
        parts = line.split('|')
        if len(parts) >= 5:
            tasks.append({
                'task_id': parts[0],
                'event_count': int(parts[1]),
                'duration': float(parts[4]) if parts[4] else 0
            })
    return tasks

def get_library_names():
    """获取库名称映射"""
    sql = "SELECT id, name FROM libraries"
    lines = run_sql(sql)
    mapping = {}
    for line in lines:
        parts = line.split('|')
        if len(parts) >= 2:
            mapping[int(parts[0])] = parts[1]
    return mapping

def get_task_details(task_id):
    """获取任务详细信息"""
    sql = f"""
    SELECT 
        al.library_id,
        al.event,
        al.created_at
    FROM activity_logs al
    WHERE al.task_id = '{task_id}'
    ORDER BY al.created_at
    """
    lines = run_sql(sql)
    events = []
    for line in lines:
        parts = line.split('|')
        if len(parts) >= 3:
            events.append({
                'library_id': int(parts[0]) if parts[0] else 0,
                'event': parts[1],
                'timestamp': parts[2]
            })
    return events

# 图表 1: 响应时间分布图 (Top 10 最快 + Top 10 最慢)
def plot_response_time_distribution(tasks):
    """绘制响应时间分布图"""
    # 过滤掉耗时为 0 的任务
    valid_tasks = [t for t in tasks if t['duration'] > 0]
    
    # 排序
    sorted_tasks = sorted(valid_tasks, key=lambda x: x['duration'])
    
    # 取前 10 和后 10
    top_10_fast = sorted_tasks[:10]
    top_10_slow = sorted_tasks[-10:]
    
    # 合并
    selected_tasks = top_10_fast + top_10_slow
    
    # 准备数据
    labels = [f"Task {i+1}" for i in range(len(selected_tasks))]
    durations = [t['duration'] for t in selected_tasks]
    colors = ['#2ecc71'] * 10 + ['#e74c3c'] * 10  # 绿色表示快，红色表示慢
    
    # 绘图
    fig, ax = plt.subplots(figsize=(14, 6))
    bars = ax.bar(labels, durations, color=colors, alpha=0.7, edgecolor='black')
    
    # 添加数值标签
    for i, (bar, duration) in enumerate(zip(bars, durations)):
        height = bar.get_height()
        if duration < 60:
            label = f'{duration:.1f}s'
        elif duration < 3600:
            label = f'{duration/60:.1f}m'
        else:
            label = f'{duration/3600:.1f}h'
        ax.text(bar.get_x() + bar.get_width()/2., height,
                label, ha='center', va='bottom', fontsize=8)
    
    ax.set_xlabel('Task Number', fontsize=12, fontweight='bold')
    ax.set_ylabel('Duration (seconds, log scale)', fontsize=12, fontweight='bold')
    ax.set_title('Response Time Distribution (Top 10 Fastest + Top 10 Slowest)', fontsize=14, fontweight='bold')
    ax.set_yscale('log')  # 使用对数刻度
    ax.grid(axis='y', alpha=0.3, which='both')
    
    # 添加图例
    from matplotlib.patches import Patch
    legend_elements = [
        Patch(facecolor='#2ecc71', alpha=0.7, label='Fastest 10 Tasks'),
        Patch(facecolor='#e74c3c', alpha=0.7, label='Slowest 10 Tasks')
    ]
    ax.legend(handles=legend_elements, loc='upper left')
    
    plt.xticks(rotation=45, ha='right')
    plt.tight_layout()
    plt.savefig('imgs/response_time_distribution.png', dpi=300, bbox_inches='tight')
    plt.close()
    print("✅ 生成: response_time_distribution.png")

# 图表 2: 性能瓶颈分析图 (饼图)
def plot_performance_bottleneck(tasks):
    """绘制性能瓶颈分析图"""
    # 使用与 analyze_performance.py 相同的统计逻辑
    from collections import defaultdict
    from datetime import datetime
    
    all_phases = defaultdict(float)
    
    for task in tasks:
        events = get_task_details(task['task_id'])
        if not events:
            continue
        
        # 使用与 analyze_performance.py 完全相同的逻辑
        prev_time = None
        prev_event_type = None
        
        for event in events:
            event_type = event['event']
            try:
                timestamp = datetime.fromisoformat(event['timestamp'].replace('+08', '+08:00'))
            except:
                continue
            
            if prev_time:
                duration = (timestamp - prev_time).total_seconds()
            else:
                duration = 0
            
            # 修正阶段归属：document.complete 的耗时实际上是 Embedding 生成的时间
            actual_event_type = event_type
            if event_type == 'document.complete' and prev_event_type == 'document.embed':
                actual_event_type = 'document.embed.actual'
            
            # 累加耗时
            all_phases[actual_event_type] += duration
            
            prev_time = timestamp
            prev_event_type = event_type
    
    # 合并 document.embed 和 document.embed.actual
    if 'document.embed.actual' in all_phases:
        all_phases['document.embed'] = all_phases.get('document.embed', 0) + all_phases['document.embed.actual']
        del all_phases['document.embed.actual']
    
    # 映射到友好的名称并计算百分比
    total = sum(all_phases.values())
    if total == 0:
        print("⚠️  没有足够的数据生成性能瓶颈分析图")
        return
    
    name_mapping = {
        'document.embed': 'Embedding Generation',
        'document.parse': 'Document Parsing',
        'document.chunk': 'Chunking',
        'github.import.download': 'GitHub Download',
        'github.import.start': 'GitHub Start'
    }
    
    stages = {}
    others = 0
    for event_type, duration in all_phases.items():
        if duration > 0:
            friendly_name = name_mapping.get(event_type, None)
            if friendly_name:
                stages[friendly_name] = round(duration / total * 100, 1)
            else:
                others += duration
    
    if others > 0:
        stages['Others'] = round(others / total * 100, 1)
    
    # 准备数据
    labels = list(stages.keys())
    sizes = list(stages.values())
    colors = ['#e74c3c', '#3498db', '#f39c12', '#2ecc71', '#95a5a6'][:len(sizes)]
    # 突出显示最大的部分
    max_idx = sizes.index(max(sizes)) if sizes else 0
    explode = tuple(0.1 if i == max_idx else 0 for i in range(len(sizes)))
    
    # 绘图
    fig, ax = plt.subplots(figsize=(10, 8))
    wedges, texts, autotexts = ax.pie(sizes, explode=explode, labels=labels, colors=colors,
                                        autopct='%1.1f%%', startangle=90, textprops={'fontsize': 11})
    
    # 美化百分比文本
    for autotext in autotexts:
        autotext.set_color('white')
        autotext.set_fontweight('bold')
        autotext.set_fontsize(12)
    
    ax.set_title('Performance Bottleneck Analysis - Time Distribution by Stage', fontsize=14, fontweight='bold', pad=20)
    
    # 添加图例
    ax.legend(wedges, [f'{label}: {size}%' for label, size in zip(labels, sizes)],
              title="Stage",
              loc="center left",
              bbox_to_anchor=(1, 0, 0.5, 1),
              fontsize=10)
    
    plt.tight_layout()
    plt.savefig('imgs/performance_bottleneck.png', dpi=300, bbox_inches='tight')
    plt.close()
    print("✅ 生成: performance_bottleneck.png")

# 图表 3: 吞吐量趋势图 (文档数量 vs 耗时)
def plot_throughput_trend(tasks):
    """绘制吞吐量趋势图"""
    # 获取每个任务的文档数量
    lib_names = get_library_names()
    task_data = []
    
    for task in tasks:
        if task['duration'] == 0:
            continue
        
        events = get_task_details(task['task_id'])
        if not events:
            continue
        
        # 统计文档数量 (document.parse 事件数)
        doc_count = sum(1 for e in events if e['event'] == 'document.parse')
        
        if doc_count > 0:
            lib_id = events[0]['library_id'] if events else 0
            task_data.append({
                'doc_count': doc_count,
                'duration': task['duration'],
                'library_id': lib_id,
                'task_id': task['task_id']
            })
            # 调试信息
            if lib_id == 17:
                print(f"[DEBUG] Task {task['task_id']}: lib_id={lib_id}, doc_count={doc_count}, duration={task['duration']:.2f}s")
    
    if not task_data:
        print("⚠️  没有足够的数据生成吞吐量趋势图")
        return
    
    # 准备数据
    doc_counts = [t['doc_count'] for t in task_data]
    durations = [t['duration'] for t in task_data]
    
    # 绘图
    fig, ax = plt.subplots(figsize=(12, 7))
    
    # 散点图
    scatter = ax.scatter(doc_counts, durations, alpha=0.6, s=100, c=durations, 
                         cmap='YlOrRd', edgecolors='black', linewidth=0.5)
    
    # 添加趋势线（在对数空间中生成点）
    if len(doc_counts) > 1:
        # 在对数空间中拟合
        log_doc_counts = np.log10(doc_counts)
        log_durations = np.log10(durations)
        z = np.polyfit(log_doc_counts, log_durations, 1)
        
        # 在对数空间中生成趋势线点
        x_trend_log = np.linspace(min(log_doc_counts), max(log_doc_counts), 100)
        y_trend_log = z[0] * x_trend_log + z[1]
        
        # 转换回线性空间
        x_trend = 10 ** x_trend_log
        y_trend = 10 ** y_trend_log
        
        ax.plot(x_trend, y_trend, "r--", alpha=0.8, linewidth=2, 
                label=f'Power Law: y={10**z[1]:.2f}x^{z[0]:.2f}')
    
    ax.set_xlabel('Document Count (log scale)', fontsize=12, fontweight='bold')
    ax.set_ylabel('Total Duration (seconds, log scale)', fontsize=12, fontweight='bold')
    ax.set_title('Throughput Trend - Document Count vs Duration', fontsize=14, fontweight='bold')
    ax.set_xscale('log')
    ax.set_yscale('log')
    ax.grid(True, alpha=0.3, which='both')
    ax.legend(fontsize=10)
    
    # 添加颜色条
    cbar = plt.colorbar(scatter, ax=ax)
    cbar.set_label('Duration (s)', fontsize=10)
    
    plt.tight_layout()
    plt.savefig('imgs/throughput_trend.png', dpi=300, bbox_inches='tight')
    plt.close()
    print("✅ 生成: throughput_trend.png")

# 图表 4: 负载与延迟关系图 (块数 vs 耗时)
def plot_load_latency_relationship(tasks):
    """绘制负载与延迟关系图"""
    # 获取每个任务的块数
    task_data = []
    
    for task in tasks:
        if task['duration'] == 0:
            continue
        
        events = get_task_details(task['task_id'])
        if not events:
            continue
        
        # 统计块数 (document.chunk 事件数)
        chunk_count = sum(1 for e in events if e['event'] == 'document.chunk')
        
        if chunk_count > 0:
            task_data.append({
                'chunk_count': chunk_count,
                'duration': task['duration']
            })
    
    if not task_data:
        print("⚠️  没有足够的数据生成负载与延迟关系图")
        return
    
    # 准备数据
    chunk_counts = [t['chunk_count'] for t in task_data]
    durations = [t['duration'] for t in task_data]
    
    # 绘图
    fig, ax = plt.subplots(figsize=(12, 7))
    
    # 散点图
    scatter = ax.scatter(chunk_counts, durations, alpha=0.6, s=100, c=chunk_counts, 
                         cmap='viridis', edgecolors='black', linewidth=0.5)
    
    # 添加趋势线（在对数空间中生成点）
    if len(chunk_counts) > 1:
        # 在对数空间中拟合
        log_chunk_counts = np.log10(chunk_counts)
        log_durations = np.log10(durations)
        z = np.polyfit(log_chunk_counts, log_durations, 1)
        
        # 在对数空间中生成趋势线点
        x_trend_log = np.linspace(min(log_chunk_counts), max(log_chunk_counts), 100)
        y_trend_log = z[0] * x_trend_log + z[1]
        
        # 转换回线性空间
        x_trend = 10 ** x_trend_log
        y_trend = 10 ** y_trend_log
        
        ax.plot(x_trend, y_trend, "r--", alpha=0.8, linewidth=2, 
                label=f'Power Law: y={10**z[1]:.2f}x^{z[0]:.2f}')
    
    ax.set_xlabel('Chunk Count (log scale)', fontsize=12, fontweight='bold')
    ax.set_ylabel('Total Duration (seconds, log scale)', fontsize=12, fontweight='bold')
    ax.set_title('Load vs Latency - Chunk Count vs Duration', fontsize=14, fontweight='bold')
    ax.set_xscale('log')
    ax.set_yscale('log')
    ax.grid(True, alpha=0.3, which='both')
    ax.legend(fontsize=10)
    
    # 添加颜色条
    cbar = plt.colorbar(scatter, ax=ax)
    cbar.set_label('Chunks', fontsize=10)
    
    plt.tight_layout()
    plt.savefig('imgs/load_latency_relationship.png', dpi=300, bbox_inches='tight')
    plt.close()
    print("✅ 生成: load_latency_relationship.png")

def main():
    print("正在获取任务数据...")
    tasks = get_task_summary()
    print(f"找到 {len(tasks)} 个任务\n")
    
    print("正在生成图表...")
    print("-" * 50)
    
    plot_response_time_distribution(tasks)
    plot_performance_bottleneck(tasks)
    plot_throughput_trend(tasks)
    plot_load_latency_relationship(tasks)
    
    print("-" * 50)
    print("\n✅ 所有图表生成完成！")
    print("图表保存位置: imgs/")

if __name__ == '__main__':
    main()
