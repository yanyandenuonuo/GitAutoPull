echo "Start auto pull:"
$GODIR		= "E:/Project/Go/"		# GO项目路径
$PHPDIR		= "E:/Project/PHP/"		# PHP项目路径
$JAVADIR 	= "E:/Project/Java/"	# Java项目路径
$HTMLDIR 	= "E:/Project/HTML/"	# HTML项目路径
$NGINXDIR 	= "E:/Project/Nginx/"	# Nginx项目路径
$CPlusDIR 	= "E:/Project/C++/"		# C++项目路径
$SQLDIR 	= "E:/Project/SQL/"		# SQL项目路径

function gitPull($path) {
	cd $path
	$dirList = get-childitem
	foreach($dir in $dirList) {
		cd $dir
		dir
		git stash
		git pull
		git stash pop
		cd ..
	}
}

gitPull($GODIR)	    # 遍历GO目录下文件
gitPull($PHPDIR)	# 遍历PHP目录下文件
gitPull($JAVADIR)	# 遍历Java目录下文件
gitPull($HTMLDIR)	# 遍历HTML目录下文件
gitPull($NGINXDIR)	# 遍历NGINX目录下文件
gitPull($CPlusDIR)	# 遍历C++目录下文件
gitPull($SQLDIR)	# 遍历SQL目录下文件

echo "End pull success!"