#!/bin/bash

# ===================== 配置参数（根据需求修改）=====================
START_DATE="2026-01-01"  # 提交起始日期（格式：YYYY-MM-DD）
END_DATE="2026-12-12"    # 提交结束日期（格式：YYYY-MM-DD）
COMMITS_PER_DAY=3        # 每天生成的提交次数
TEMP_FILE=".git-demo-commit.tmp"  # 临时提交文件，脚本结束后可删除
# =================================================================

# 检查是否在Git仓库目录下
if [ ! -d ".git" ]; then
    echo "错误：当前目录不是Git仓库，请先执行 git init 初始化仓库！"
    exit 1
fi

# BSD date 转换日期格式（兼容Mac），将 YYYY-MM-DD 转为时间戳用于比较
current_ts=$(date -j -f "%Y-%m-%d" "$START_DATE" "+%s")
end_ts=$(date -j -f "%Y-%m-%d" "$END_DATE" "+%s")

# 循环遍历每一天，生成提交记录
while [ "$current_ts" -le "$end_ts" ]; do
    # 将时间戳转回 YYYY-MM-DD 格式
    current_date=$(date -j -f "%s" "$current_ts" "+%Y-%m-%d")
    echo "正在生成 $current_date 的提交记录（共 $COMMITS_PER_DAY 次）..."
    
    # 当天内生成指定次数的提交
    for ((i=1; i<=COMMITS_PER_DAY; i++)); do
        # 1. 修改临时文件（确保有文件变更，Git才能提交）
        echo "提交记录：$current_date - 第 $i 次提交" >> "$TEMP_FILE"
        
        # 2. 添加到Git暂存区
        git add "$TEMP_FILE"
        
        # 3. 提交（指定提交日期，随机生成当天8-10点的时间，兼容Mac）
        random_hour=$((8 + RANDOM % 2))
        random_min=$((RANDOM % 60))
        random_sec=$((RANDOM % 60))
        commit_date="$current_date $random_hour:$random_min:$random_sec"
        
        GIT_AUTHOR_DATE="$commit_date" \
        GIT_COMMITTER_DATE="$commit_date" \
        git commit -m "feat: 模拟 $current_date 第 $i 次提交 [自动生成]"
        
        # 随机休眠（可选，模拟真实提交间隔，注释后可快速生成）
        sleep 1
    done
    
    # 推进到下一天（BSD date 用 -v +1d 增加一天，再转为时间戳）
    next_date=$(date -j -f "%Y-%m-%d" "$current_date" -v +1d "+%Y-%m-%d")
    current_ts=$(date -j -f "%Y-%m-%d" "$next_date" "+%s")
done

# 可选：删除临时提交文件（如需保留提交痕迹，可注释该行）
rm -f "$TEMP_FILE"

echo "============================================="
echo "模拟Git提交记录生成完成！可执行 git log --graph --pretty=oneline 查看结果"