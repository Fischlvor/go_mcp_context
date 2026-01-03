#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
生成文档检索性能图表
从 mcpcalllog 表查询 search-libraries 和 get-library-docs 的耗时数据
"""

import subprocess
import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns
import numpy as np
from datetime import datetime
import sys
import os

# 添加父目录到路径以导入 db_config
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', 'docs_update'))
from db_config import get_sql_command

# 设置中文字体
plt.rcParams['font.sans-serif'] = ['DejaVu Sans', 'Arial Unicode MS', 'SimHei']
plt.rcParams['axes.unicode_minus'] = False

def run_sql(sql):
    """执行 SQL 查询并返回结果"""
    cmd = get_sql_command(sql)
    result = subprocess.run(cmd, shell=True, capture_output=True, text=True)
    if result.returncode != 0:
        print(f"SQL 执行失败: {result.stderr}", file=sys.stderr)
        return []
    return [line.strip() for line in result.stdout.strip().split('\n') if line.strip()]

def fetch_data():
    """从数据库获取数据"""
    print("查询数据...")
    
    query = """
    SELECT 
        params->'params'->>'name' as tool_name,
        library_id,
        latency_ms,
        created_at
    FROM mcp_call_logs
    WHERE func_name = 'tools_call'
        AND created_at >= '2026-01-03 00:00:00'
        AND (
            (params->'params'->>'name' = 'search-libraries' AND latency_ms >= 1)
            OR
            (params->'params'->>'name' = 'get-library-docs' AND latency_ms >= 50)
        )
    ORDER BY created_at DESC
    """
    
    # 执行 SQL
    lines = run_sql(query)
    
    # 解析结果
    data = []
    for line in lines:
        if not line:
            continue
        parts = line.split('|')
        if len(parts) >= 4:
            try:
                tool_name = parts[0].strip().strip('"')
                data.append({
                    'func_name': tool_name,
                    'library_id': parts[1].strip() if parts[1].strip() else None,
                    'latency_ms': float(parts[2]),
                    'created_at': parts[3].strip()
                })
            except (ValueError, IndexError):
                continue
    
    df = pd.DataFrame(data)
    
    print(f"获取到 {len(df)} 条记录")
    if len(df) > 0:
        print(f"  search-libraries: {len(df[df['func_name'] == 'search-libraries'])} 条")
        print(f"  get-library-docs: {len(df[df['func_name'] == 'get-library-docs'])} 条")
    
    return df

def print_statistics(df):
    """打印统计信息"""
    print("\n" + "="*60)
    print("统计信息")
    print("="*60)
    
    for func_name in ['search-libraries', 'get-library-docs']:
        data = df[df['func_name'] == func_name]['latency_ms']
        if len(data) == 0:
            continue
            
        print(f"\n{func_name}:")
        print(f"  样本数: {len(data)}")
        print(f"  均值: {data.mean():.2f} ms")
        print(f"  中位数: {data.median():.2f} ms")
        print(f"  标准差: {data.std():.2f} ms")
        print(f"  最小值: {data.min():.2f} ms")
        print(f"  最大值: {data.max():.2f} ms")
        print(f"  P25: {data.quantile(0.25):.2f} ms")
        print(f"  P75: {data.quantile(0.75):.2f} ms")
        print(f"  P95: {data.quantile(0.95):.2f} ms")
        print(f"  P99: {data.quantile(0.99):.2f} ms")

def create_boxplot(df, output_path):
    """创建箱线图（单图双 Y 轴）"""
    print("\n生成箱线图...")
    
    fig, ax1 = plt.subplots(figsize=(10, 6))
    
    # 准备数据
    search_data = df[df['func_name'] == 'search-libraries']['latency_ms']
    docs_data = df[df['func_name'] == 'get-library-docs']['latency_ms']
    
    # 创建第二个 Y 轴
    ax2 = ax1.twinx()
    
    # 在 ax1 上绘制 search-libraries（左 Y 轴）
    if len(search_data) > 0:
        bp1 = ax1.boxplot([search_data], positions=[1], widths=0.6, patch_artist=True,
                          showmeans=True, meanline=True, showfliers=True,
                          medianprops={'color': 'darkred', 'linewidth': 2},
                          meanprops={'color': 'darkblue', 'linewidth': 2, 'linestyle': '--'},
                          boxprops={'facecolor': '#3498db', 'alpha': 0.6},
                          flierprops={'marker': 'o', 'markerfacecolor': '#3498db', 'markersize': 3, 'alpha': 0.5})
        
        ax1.set_ylabel('search-libraries Latency (ms)', fontsize=12, color='#3498db')
        ax1.tick_params(axis='y', labelcolor='#3498db')
        ax1.set_ylim(0, max(search_data.quantile(0.95) * 1.2, 100))
        ax1.grid(True, alpha=0.2, axis='y', linestyle='--')
    
    # 在 ax2 上绘制 get-library-docs（右 Y 轴）
    if len(docs_data) > 0:
        bp2 = ax2.boxplot([docs_data], positions=[2], widths=0.6, patch_artist=True,
                          showmeans=True, meanline=True, showfliers=True,
                          medianprops={'color': 'darkred', 'linewidth': 2},
                          meanprops={'color': 'darkblue', 'linewidth': 2, 'linestyle': '--'},
                          boxprops={'facecolor': '#2ecc71', 'alpha': 0.6},
                          flierprops={'marker': 'o', 'markerfacecolor': '#2ecc71', 'markersize': 3, 'alpha': 0.5})
        
        ax2.set_ylabel('get-library-docs Latency (ms)', fontsize=12, color='#2ecc71')
        ax2.tick_params(axis='y', labelcolor='#2ecc71')
        ax2.set_ylim(0, max(docs_data.quantile(0.95) * 1.2, 2500))
        ax2.grid(True, alpha=0.2, axis='y', linestyle=':', color='#2ecc71')
    
    # 设置 X 轴
    ax1.set_xticks([1, 2])
    ax1.set_xticklabels(['search-libraries', 'get-library-docs'])
    ax1.set_xlim(0.5, 2.5)
    
    # 标题
    ax1.set_title('MCP Function Call Latency - Box Plot (Dual Y-Axis)', 
                  fontsize=14, fontweight='bold', pad=20)
    
    # 添加图例
    from matplotlib.lines import Line2D
    legend_elements = [
        Line2D([0], [0], color='darkred', linewidth=2, label='Median'),
        Line2D([0], [0], color='darkblue', linewidth=2, linestyle='--', label='Mean'),
        Line2D([0], [0], marker='s', color='w', markerfacecolor='#3498db', 
               markersize=10, alpha=0.6, label='search-libraries'),
        Line2D([0], [0], marker='s', color='w', markerfacecolor='#2ecc71', 
               markersize=10, alpha=0.6, label='get-library-docs')
    ]
    ax1.legend(handles=legend_elements, loc='upper center', frameon=True, fancybox=True)
    
    # 添加统计信息
    if len(search_data) > 0:
        stats1 = f"search-libraries:\nn={len(search_data)}, Median={search_data.median():.1f}ms, P95={search_data.quantile(0.95):.1f}ms"
        ax1.text(0.02, 0.98, stats1, transform=ax1.transAxes,
                fontsize=9, verticalalignment='top',
                bbox=dict(boxstyle='round', facecolor='lightblue', alpha=0.7))
    
    if len(docs_data) > 0:
        stats2 = f"get-library-docs:\nn={len(docs_data)}, Median={docs_data.median():.1f}ms, P95={docs_data.quantile(0.95):.1f}ms"
        ax1.text(0.98, 0.98, stats2, transform=ax1.transAxes,
                fontsize=9, verticalalignment='top', horizontalalignment='right',
                bbox=dict(boxstyle='round', facecolor='lightgreen', alpha=0.7))
    
    plt.tight_layout()
    plt.savefig(output_path, dpi=300, bbox_inches='tight')
    print(f"箱线图已保存: {output_path}")
    plt.close()

def create_histogram(df, output_path):
    """创建直方图"""
    print("\n生成直方图...")
    
    fig, ax = plt.subplots(figsize=(12, 6))
    
    # 准备数据
    search_data = df[df['func_name'] == 'search-libraries']['latency_ms']
    docs_data = df[df['func_name'] == 'get-library-docs']['latency_ms']
    
    # 确定合适的 bins
    max_duration = max(search_data.max() if len(search_data) > 0 else 0,
                       docs_data.max() if len(docs_data) > 0 else 0)
    bins = np.linspace(0, min(max_duration, 2000), 50)  # 限制在 3000ms 以内
    
    # 绘制直方图
    if len(search_data) > 0:
        weights_search = np.ones_like(search_data)
        ax.hist(search_data, bins=bins, weights=weights_search, alpha=0.6, 
                label='search-libraries', color='#3498db', edgecolor='black')
    
    if len(docs_data) > 0:
        ax.hist(docs_data, bins=bins, alpha=0.6, label='get-library-docs',
                color='#2ecc71', edgecolor='black')
    
    ax.set_xlabel('Latency (ms)', fontsize=12)
    ax.set_ylabel('Count', fontsize=12)
    ax.set_title('MCP Function Call Latency Distribution', fontsize=14, fontweight='bold')
    ax.legend(loc='upper right', fontsize=10)
    ax.grid(True, alpha=0.3, axis='y')
    
    plt.tight_layout()
    plt.savefig(output_path, dpi=300, bbox_inches='tight')
    print(f"直方图已保存: {output_path}")
    plt.close()

def create_combined_chart(df, output_path):
    """创建组合图表（箱线图 + 统计信息）"""
    print("\n生成组合图表...")
    
    fig, (ax1, ax2) = plt.subplots(1, 2, figsize=(16, 6))
    
    # 左侧：箱线图
    data_to_plot = []
    labels = []
    stats_text = []
    
    for func_name in ['search-libraries', 'get-library-docs']:
        data = df[df['func_name'] == func_name]['latency_ms']
        if len(data) > 0:
            data_to_plot.append(data)
            labels.append(func_name.replace('-', '\n'))
            
            # 准备统计文本
            stats = f"n={len(data)}\n"
            stats += f"Mean: {data.mean():.1f}ms\n"
            stats += f"Median: {data.median():.1f}ms\n"
            stats += f"P95: {data.quantile(0.95):.1f}ms"
            stats_text.append(stats)
    
    bp = ax1.boxplot(data_to_plot, labels=labels, patch_artist=True,
                      showmeans=True, meanline=True,
                      medianprops={'color': 'red', 'linewidth': 2},
                      meanprops={'color': 'blue', 'linewidth': 2, 'linestyle': '--'})
    
    colors = ['#3498db', '#2ecc71']
    for patch, color in zip(bp['boxes'], colors):
        patch.set_facecolor(color)
        patch.set_alpha(0.6)
    
    ax1.set_ylabel('Duration (ms)', fontsize=12)
    ax1.set_title('Box Plot with Statistics', fontsize=14, fontweight='bold')
    ax1.grid(True, alpha=0.3, axis='y')
    
    # 添加统计文本
    for i, (label, stats) in enumerate(zip(labels, stats_text)):
        ax1.text(i+1, ax1.get_ylim()[1] * 0.95, stats,
                ha='center', va='top', fontsize=9,
                bbox=dict(boxstyle='round', facecolor='white', alpha=0.8))
    
    # 右侧：直方图
    search_data = df[df['func_name'] == 'search-libraries']['latency_ms']
    docs_data = df[df['func_name'] == 'get-library-docs']['latency_ms']
    
    max_duration = max(search_data.max() if len(search_data) > 0 else 0,
                       docs_data.max() if len(docs_data) > 0 else 0)
    bins = np.linspace(0, min(max_duration, 2000), 40)
    
    if len(search_data) > 0:
        ax2.hist(search_data, bins=bins, alpha=0.6, label='search-libraries',
                color='#3498db', edgecolor='black')
    
    if len(docs_data) > 0:
        ax2.hist(docs_data, bins=bins, alpha=0.6, label='get-library-docs',
                color='#2ecc71', edgecolor='black')
    
    ax2.set_xlabel('Duration (ms)', fontsize=12)
    ax2.set_ylabel('Count', fontsize=12)
    ax2.set_title('Duration Distribution', fontsize=14, fontweight='bold')
    ax2.legend(loc='upper right', fontsize=10)
    ax2.grid(True, alpha=0.3, axis='y')
    
    plt.tight_layout()
    plt.savefig(output_path, dpi=300, bbox_inches='tight')
    print(f"组合图表已保存: {output_path}")
    plt.close()

def main():
    """主函数"""
    print("="*60)
    print("MCP 文档检索性能图表生成工具")
    print("="*60)
    
    # 获取数据
    df = fetch_data()
    
    if len(df) == 0:
        print("错误: 没有获取到数据")
        return
    
    # 打印统计信息
    print_statistics(df)
    
    # 创建输出目录
    import os
    output_dir = os.path.dirname(os.path.abspath(__file__))
    imgs_dir = os.path.join(output_dir, 'imgs')
    os.makedirs(imgs_dir, exist_ok=True)
    
    # 生成图表
    create_boxplot(df, os.path.join(imgs_dir, 'mcp_duration_boxplot.png'))
    create_histogram(df, os.path.join(imgs_dir, 'mcp_duration_histogram.png'))
    
    print("\n" + "="*60)
    print("所有图表生成完成！")
    print("="*60)

if __name__ == '__main__':
    main()
