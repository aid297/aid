package main

import (
	"flag"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

var (
	branchSrc string
	branchDst string
	comment   *string
	args      []string
	branch    *string
	tag       *string
)

func main() {
	branch = flag.String("branch", "", "分支")
	tag = flag.String("tag", "", "标签")
	comment = flag.String("comment", "", "合并注释")

	flag.Parse()
	args = flag.Args()

	if len(args) < 2 {
		log.Fatalln("参数不足，至少需要一个参数")
	}

	switch args[0] {
	case "push":
		err := exec.Command("git", "add", "--all").Run()
		if err != nil {
			log.Fatalf("执行 git add 失败：%v", err)
		}

		if *comment == "" {
			log.Fatalln("提交注释不能为空")
		}
		err = exec.Command("git", "commit", "-m", *comment).Run()
		if err != nil {
			log.Fatalf("执行 git commit 失败：%v", err)
		}

		if *branch == "" {
			err = exec.Command("git", "pull").Run()
			if err != nil {
				log.Fatalf("执行 git pull 失败：%v", err)
			}

			err = exec.Command("git", "push").Run()
			if err != nil {
				log.Fatalf("执行 git push 失败：%v", err)
			}
		} else {
			err = exec.Command("git", "pull").Run()
			if err != nil {
				log.Fatalf("执行 git pull 失败：%v", err)
			}

			err = exec.Command("git", "push", "-u", "origin", *branch).Run()
			if err != nil {
				log.Fatalf("执行 git push 失败：%v", err)
			}
		}

	case "merge":
		if *branch == "" {
			log.Fatalln("分支不能为空")
		}

		branchParts := strings.Split(*branch, " -> ")
		if len(branchParts) != 2 {
			log.Fatalln("分支格式错误，需要为 src -> dst")
		}
		branchSrc = branchParts[0]
		branchDst = branchParts[1]

		err := exec.Command("git", "checkout", branchDst).Run()
		if err != nil {
			log.Fatalf("执行 git checkout %s 失败：%v", branchDst, err)
		}

		err = exec.Command("git", "fetch", "origin").Run()
		if err != nil {
			log.Fatalf("执行 git fetch 失败：%v", err)
		}

		err = exec.Command("git", "pull", "origin", branchDst).Run()
		if err != nil {
			log.Fatalf("执行 git pull 失败：%v", err)
		}

		err = exec.Command("git", "merge", fmt.Sprintf("origin/%s", branchSrc)).Run()
		if err != nil {
			log.Fatalf("执行 git merge 失败：%v", err)
		}

		err = exec.Command("git", "merge", "-u", "origin", branchDst).Run()
		if err != nil {
			log.Fatalf("执行 git merge -u origin/%s 失败：%v", branchDst, err)
		}

		err = exec.Command("git", "checkout", branchSrc).Run()
		if err != nil {
			log.Fatalf("执行 git checkout %s 失败：%v", branchSrc, err)
		}
	case "last-tag":
		cmd := exec.Command("git", "describe", "--tags", "$(git rev-list --tags --max-count=1)")
		out, err := cmd.Output()
		if err != nil {
			log.Fatalf("执行 git describe 失败：%v", err)
		}
		fmt.Println("最新标签:", string(out))
	case "tag-push":
		if *tag == "" {
			log.Fatalln("标签不能为空")
		}

		err := exec.Command("git", "tag", *tag).Run()
		if err != nil {
			log.Fatalf("执行 git tag %s 失败：%v", *tag, err)
		}

		err = exec.Command("git", "push", "-u", "origin", *tag).Run()
		if err != nil {
			log.Fatalf("执行 git push 失败：%v", err)
		}
	}
}
