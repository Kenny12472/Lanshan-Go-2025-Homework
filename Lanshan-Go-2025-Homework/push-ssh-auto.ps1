# -----------------------------
# 一键使用 SSH 推送 Go 项目到 GitHub（自动检测）
# -----------------------------

# 配置
$GitHubUser = "Kenny12472"          # GitHub 用户名
$RepoName = "Lanshan-Go-2025-Homework"  # 仓库名
$Email = "你的邮箱@example.com"      # GitHub 邮箱

# 当前目录
$projectDir = Get-Location
Write-Host "当前目录：" $projectDir

# SSH 密钥路径
$sshPath = "$env:USERPROFILE\.ssh\id_rsa"

# 1️⃣ 检查 SSH 密钥
if (-Not (Test-Path $sshPath)) {
    Write-Host "没有检测到 SSH 密钥，正在生成..."
    ssh-keygen -t rsa -b 4096 -C $Email -f $sshPath -N ""
} else {
    Write-Host "已存在 SSH 密钥，跳过生成。"
}

# 2️⃣ 测试 SSH 连接 GitHub
Write-Host "`n测试 SSH 连接 GitHub..."
$sshTest = ssh -o StrictHostKeyChecking=no -T git@github.com 2>&1

if ($sshTest -match "successfully authenticated") {
    Write-Host "SSH 已经可以连接 GitHub ✅"
} else {
    Write-Host "`nSSH 还没有添加到 GitHub！"
    Write-Host "请将以下公钥复制到 GitHub → Settings → SSH and GPG keys → New SSH key:`n"
    Get-Content "$sshPath.pub" | Write-Host
    Write-Host "`n完成后按回车继续..."
    Read-Host
}

# 3️⃣ 设置远程仓库为 SSH
$sshUrl = "git@github.com:$GitHubUser/$RepoName.git"
Write-Host "`n设置远程仓库为 SSH 地址： $sshUrl"
git remote set-url origin $sshUrl

# 4️⃣ 推送 main 分支
Write-Host "`n正在推送 main 分支到 GitHub..."
git push -u origin main

Write-Host "`n操作完成 ✅"
