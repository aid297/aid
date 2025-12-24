operator="$1"

if [ "$operator" = "push-tag" ]; then
	tag="$2"
	git tag "$tag" && git push origin "$tag"
elif [ "$operator" = "last-tag" ]; then
	git describe --tags $(git rev-list --tags --max-count=1)
elif [ "$operator" = "to" ]; then
	git checkout $1
elif [ "$operator" = "help" ]; then
	echo "    1、基本用法：source git.sh 执行子程序名称 [参数1, 参数2...]"
	echo "    2、子程序名称为push-tag  ：创建并推送当前分支内容到tag（参数1表示tag名称）"
	echo "    3、子程序名称为last-tag  ：获取当前最后一个tag名称"
	echo "    4、其他子程序名称        ：分为三种情况"
	echo "        4.1、source git.sh : 提交说明          ：代表当前分支内容只提交不推送"
	echo "        4.2、source git.sh dev 提交说明        ：代表当前内容推送到dev分支"
	echo "        4.3、source git.sh dev:master 提交说明 ：代表先推送到dev分支然后合并到master再推送到master"
else
	IFS=':' read -r src_branch dst_branch <<<"$operator"
	commit="$2"
	tag="$3"

	if [ -n "$commit" ]; then
		git add --all && git commit -m"$commit"
	fi

	if [ -n "$src_branch" ]; then
		git push origin $src_branch

		if [ -n "$dst_branch" ]; then
			git checkout $dst_branch &&
				git merge $src_branch &&
				git pull origin $dst_branch &&
				git push origin $dst_branch &&
				git checkout $src_branch
		fi
	fi

	if [ -n "$tag" ]; then
		git tag "$tag" && git push origin "$tag"
	fi
fi
