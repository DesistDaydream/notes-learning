# 创建符号链接的函数
function Create-SymbolicLinks {
    param (
        [Parameter(Mandatory=$true)]
        [string]$SourceDir,  # B 目录（原始文件所在目录）

        [Parameter(Mandatory=$true)]
        [string]$TargetDir,  # A 目录（创建链接的目录）

        [string]$Filter = "*"  # 文件过滤器，默认所有文件
    )

    # 检查源目录是否存在
    if (-not (Test-Path $SourceDir)) {
        Write-Error "源目录不存在: $SourceDir"
        return
    }

    # 如果目标目录不存在，创建它
    if (-not (Test-Path $TargetDir)) {
        New-Item -ItemType Directory -Path $TargetDir -Force | Out-Null
        Write-Host "已创建目标目录: $TargetDir"
    }

    # 获取源目录中的文件
    $files = Get-ChildItem -Path $SourceDir -Filter $Filter -File

    if ($files.Count -eq 0) {
        Write-Warning "未找到匹配的文件"
        return
    }

    Write-Host "找到 $($files.Count) 个文件"

    # 为每个文件创建符号链接
    foreach ($file in $files) {
        $linkPath = Join-Path $TargetDir $file.Name

        # 如果链接已存在，跳过或删除
        if (Test-Path $linkPath) {
            Write-Warning "链接已存在，跳过: $($file.Name)"
            continue
        }

        try {
            # 创建符号链接
            New-Item -ItemType SymbolicLink -Path $linkPath -Target $file.FullName | Out-Null
            Write-Host "✓ 已创建链接: $($file.Name)"
        }
        catch {
            Write-Error "创建链接失败 $($file.Name): $_"
        }
    }
}


Create-SymbolicLinks `
-SourceDir "D:\Projects\DesistDaydream\notes-learning\content\zh-cn\.obsidian\plugins\manual-sorting" `
-TargetDir "D:\Projects\DesistDaydream\notes-pastime\.obsidian\plugins\manual-sorting"

Create-SymbolicLinks `
-SourceDir "D:\Projects\DesistDaydream\notes-learning\content\zh-cn\.obsidian\plugins\manual-sorting" `
-TargetDir "D:\Projects\DesistDaydream\notes-science\.obsidian\plugins\manual-sorting"

Create-SymbolicLinks `
-SourceDir "D:\Projects\DesistDaydream\notes-learning\content\zh-cn\.obsidian\plugins\manual-sorting" `
-TargetDir "D:\Projects\DesistDaydream\notes-haohan\.obsidian\plugins\manual-sorting"

# Create-SymbolicLinks -SourceDir "C:\B目录" -TargetDir "C:\A目录" -Filter "*.txt"