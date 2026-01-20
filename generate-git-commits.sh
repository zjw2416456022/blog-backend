#!/bin/bash

# ===================== 配置参数（随机深度+随机断更，可按需调整）=====================
START_DATE="2025-10-01"  # 提交起始日期（YYYY-MM-DD）
END_DATE="2025-10-31"    # 提交结束日期（2099年）
TEMP_FILE=".git-demo-commit.tmp"  # 临时提交文件
GITHUB_BRANCH="main"     # 你的GitHub仓库默认分支（main/master）
GITHUB_REPO="git@github.com:zjw2416456022/blog-backend.git"  # 你的GitHub仓库地址

# 1. 随机断更配置：设置当天无提交的概率（数值越小，断更概率越低，更贴近真实）
# 取值范围：1-10（推荐3，即30%概率当天无提交，随机隔1-2天断更）
SKIP_PROBABILITY=3

# 2. 颜色深度对应提交次数（匹配GitHub贡献图，随机选择其中一个区间）
# 格式：[最小次数, 最大次数]，对应不同绿色深度
COLOR_DEPTH_CONFIGS=(
    "1,1"   # 最浅绿色（1次提交）
    "2,3"   # 中浅绿色（2-3次提交）
    "4,5"   # 中深绿色（4-5次提交）
    "6,8"   # 最深绿色（6-8次提交）
)
# =================================================================

# 关键：强制英文区域，彻底避免中文日期格式干扰
export LC_ALL=C
export LANG=C

# 1. 检查并初始化Git仓库（关联GitHub远程仓库）
if [ ! -d ".git" ]; then
    echo "当前目录不是Git仓库，正在初始化并关联GitHub远程仓库..."
    git init
    git remote add origin "$GITHUB_REPO"
fi

# 2. 验证远程仓库关联是否正常
git fetch origin "$GITHUB_BRANCH" > /dev/null 2>&1
if [ $? -ne 0 ]; then
    echo "警告：远程仓库关联失败，请检查GITHUB_REPO地址是否正确！"
    echo "若已手动关联，可忽略此警告，继续生成本地提交..."
fi

# 3. 预处理临时文件（清空已有内容，确保每次提交有唯一变更）
> "$TEMP_FILE"

# 4. 日期转时间戳（仅执行2次，后续直接通过时间戳加减推进，规避格式解析）
start_ts=$(date -j -f "%Y-%m-%d" "$START_DATE" "+%s")
end_ts=$(date -j -f "%Y-%m-%d" "$END_DATE" "+%s")

# 验证时间戳转换是否成功
if [ -z "$start_ts" ] || [ -z "$end_ts" ] || [ "$start_ts" -gt "$end_ts" ]; then
    echo "错误：日期格式转换失败或起始日期晚于结束日期！"
    echo "请检查START_DATE和END_DATE格式是否为YYYY-MM-DD"
    exit 1
fi

# 5. 循环生成提交记录（核心：时间戳+86400秒推进一天，无格式解析问题）
current_ts="$start_ts"
while [ "$current_ts" -le "$end_ts" ]; do
    # 将当前时间戳转回YYYY-MM-DD格式（仅用于提交信息和日志）
    current_date=$(date -j -f "%s" "$current_ts" "+%Y-%m-%d")
    echo "============================================="
    
    # 6. 随机判断：当天是否跳过提交（实现隔1-2天无提交）
    skip_today=$((RANDOM % 10))  # 生成0-9的随机数
    if [ "$skip_today" -lt "$SKIP_PROBABILITY" ]; then
        echo "$current_date ：随机跳过提交（模拟真实断更）"
        # 直接推进到下一天，不生成任何提交
        current_ts=$((current_ts + 86400))
        continue
    fi

    # 7. 随机选择颜色深度（对应提交次数区间）
    # 随机获取配置索引
    depth_index=$((RANDOM % ${#COLOR_DEPTH_CONFIGS[@]}))
    depth_config=${COLOR_DEPTH_CONFIGS[$depth_index]}
    # 拆分最小/最大提交次数
    min_commits=$(echo "$depth_config" | cut -d ',' -f 1)
    max_commits=$(echo "$depth_config" | cut -d ',' -f 2)
    # 随机生成当天提交次数（在选中的颜色深度区间内）
    commits_today=$((min_commits + RANDOM % (max_commits - min_commits + 1)))

    # 8. 打印当天配置（颜色深度对应信息）
    depth_desc=""
    if [ "$min_commits" -eq 1 ]; then
        depth_desc="（最浅绿色）"
    elif [ "$min_commits" -eq 2 ]; then
        depth_desc="（中浅绿色）"
    elif [ "$min_commits" -eq 4 ]; then
        depth_desc="（中深绿色）"
    elif [ "$min_commits" -eq 6 ]; then
        depth_desc="（最深绿色）"
    fi
    echo "正在生成 $current_date 的提交记录（共 $commits_today 次 $depth_desc）"
    
    # 9. 生成当天的多次提交
    for ((i=1; i<=commits_today; i++)); do
        # 追加唯一内容到临时文件（确保Git检测到文件变更，避免重复）
        commit_unique_id="$current_date-$i-$(date +%s%N)"
        echo "提交记录：$commit_unique_id（GitHub随机颜色深度专用）" >> "$TEMP_FILE"
        
        # Git暂存+提交（指定提交日期，确保GitHub统计生效）
        git add "$TEMP_FILE"
        random_hour=$((9 + RANDOM % 9))  # 9-17点，贴近真实开发时间
        random_min=$((RANDOM % 60))
        random_sec=$((RANDOM % 60))
        commit_datetime="$current_date $random_hour:$random_min:$random_sec"
        
        # 执行提交（隐藏冗余日志，仅显示进度）
        GIT_AUTHOR_DATE="$commit_datetime" \
        GIT_COMMITTER_DATE="$commit_datetime" \
        git commit -m "feat: $current_date 第 $i 次提交 [随机深度$depth_desc]" > /dev/null 2>&1
        
        # 打印提交进度
        echo -n "已完成 $i/$commits_today 次提交..."$'\r'
    done
    echo -e "\n$current_date 提交生成完成！"
    
    # 10. 推进到下一天（纯时间戳计算，无格式解析漏洞）
    current_ts=$((current_ts + 86400))
    
    # 11. 每30天推送到GitHub一次（避免频繁推送被限制，提升效率）
    days_passed=$(( (current_ts - start_ts) / 86400 ))
    if [ $((days_passed % 30)) -eq 0 ] && [ $days_passed -ne 0 ]; then
        echo "============================================="
        echo "已生成 $days_passed 天提交，正在推送到GitHub远程仓库..."
        git push -u origin "$GITHUB_BRANCH" > /dev/null 2>&1
        if [ $? -eq 0 ]; then
            echo "推送成功！当前进度：$current_date"
        else
            echo "推送失败！请检查网络或仓库权限，可手动执行 git push origin $GITHUB_BRANCH"
        fi
    fi
done

# 12. 最后一次推送剩余所有提交到GitHub
echo "============================================="
echo "所有日期提交生成完成，正在推送最后一批记录到GitHub..."
git push -u origin "$GITHUB_BRANCH"
if [ $? -eq 0 ]; then
    echo "全部推送成功！GitHub贡献图将显示随机深度（需等待几分钟同步）"
else
    echo "最后一批推送失败，请手动执行 git push origin $GITHUB_BRANCH 完成推送"
fi

# 13. 清理临时文件
rm -f "$TEMP_FILE"

echo "============================================="
echo "操作全部完成！可访问GitHub仓库查看贡献图：https://github.com/zjw2416456022/blog-backend"