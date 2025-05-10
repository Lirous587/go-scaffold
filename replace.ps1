# 读取 vars.env 变量
$vars = @{}
Get-Content vars.env | ForEach-Object {
    if ($_ -match '^(.*?)=(.*)$') {
        $vars[$matches[1]] = $matches[2]
    }
}

# 定义要处理的目录
$inputDir = "./"

# 遍历所有文件并原地替换
Get-ChildItem -Path $inputDir -Recurse -File | ForEach-Object {
    $content = Get-Content $_.FullName -Raw
    foreach ($key in $vars.Keys) {
        $content = $content -replace "\$\{$key\}", $vars[$key]
    }
    Set-Content $_.FullName $content
}
