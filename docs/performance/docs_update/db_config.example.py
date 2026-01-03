#!/usr/bin/env python3
"""
数据库连接配置模板
复制此文件为 db_config.py 并填入实际的连接信息
"""

# SSH 连接配置
SSH_HOST = "user@your-server-ip"

# PostgreSQL 配置
POSTGRES_CONTAINER = "your_postgres_container"
POSTGRES_USER = "postgres"
POSTGRES_DB = "your_database"

def get_sql_command(sql):
    """
    生成完整的 SSH + PostgreSQL 命令
    
    Args:
        sql: SQL 查询语句
        
    Returns:
        完整的 SSH 命令字符串
    """
    return f'ssh {SSH_HOST} "docker exec {POSTGRES_CONTAINER} psql -U {POSTGRES_USER} -d {POSTGRES_DB} -t -A -F \'|\' -c \\"{sql}\\""'
