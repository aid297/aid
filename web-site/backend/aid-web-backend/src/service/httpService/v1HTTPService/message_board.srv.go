package v1HTTPService

import (
	"errors"
	"fmt"
	"os"

	"github.com/gofrs/uuid/v5"
	jsonIter "github.com/json-iterator/go"

	"github.com/aid297/aid/v2/anySlices"
	"github.com/aid297/aid/v2/filesystems"
	"github.com/aid297/aid/v2/web-site/backend/aid-web-backend/src/global"
	"github.com/aid297/aid/v2/web-site/backend/aid-web-backend/src/module/httpModule/v1HTTPModule/request"
)

// MessageBoardService 服务：留言板
type MessageBoardService struct{}

func (*MessageBoardService) getDirectionFile() (
	directionFile filesystems.Filesystem,
	directionFileSlice anySlices.AnySlicer[string],
	err error,
) {
	var (
		directionDir         = filesystems.NewDir(filesystems.Rel(global.CONFIG.MessageBoard.Dir))
		directionFileJSON    []byte
		directionFileContent []string
	)

	if !directionDir.GetExist() {
		if err = directionDir.Create().GetError(); err != nil {
			return
		}
	}

	directionFile = filesystems.NewFile(filesystems.Rel(global.CONFIG.MessageBoard.Dir, "direction.json"))

	if !directionDir.GetExist() {
		return nil, nil, nil
	}

	if !directionFile.GetExist() {
		if err = directionFile.Write([]byte{'[', ']'}, filesystems.Mode(0755), filesystems.Flag(os.O_WRONLY|os.O_TRUNC|os.O_CREATE)).GetError(); err != nil {
			return nil, nil, fmt.Errorf("创建目录文件失败：%w", err)
		}
		return directionFile, anySlices.New[string](), nil
	}

	if directionFileJSON, err = directionFile.Read(); err != nil {
		return
	}

	if err = jsonIter.Unmarshal(directionFileJSON, &directionFileContent); err != nil {
		return
	}

	directionFileSlice = anySlices.NewList(directionFileContent)

	return
}

// List 留言板服务：获取信息列表
func (*MessageBoardService) List() (messages []map[string]string, err error) {
	var (
		directionFileSlice anySlices.AnySlicer[string]
		messageFile        filesystems.Filesystem
		messageFileContent []byte
		messageContent     map[string]string
	)

	if _, directionFileSlice, err = New.MessageBoard().getDirectionFile(); err != nil {
		return
	}

	messages = make([]map[string]string, directionFileSlice.LengthNotEmpty())
	directionFileSlice.Each(func(idx int, item string) (isBreak bool) {
		messageFile = filesystems.NewFile(filesystems.Rel(global.CONFIG.MessageBoard.Dir, item))
		if !messageFile.GetExist() {
			return
		}

		messageContent = make(map[string]string)
		if messageFileContent, err = messageFile.Read(); err != nil {
			return
		}

		if err = jsonIter.Unmarshal(messageFileContent, &messageContent); err != nil {
			return
		}

		messages[idx] = messageContent

		return
	})

	return messages, nil
}

// Store 留言板服务：保存信息
func (*MessageBoardService) Store(form *request.MessageBoardStoreRequest) (err error) {
	var (
		directionFile      filesystems.Filesystem
		directionFileJSON  []byte
		directionFileSlice anySlices.AnySlicer[string]
		newUUID            = uuid.Must(uuid.NewV7())
		messageFile        filesystems.Filesystem
		messageFileContent []byte
	)

	if directionFile, directionFileSlice, err = New.MessageBoard().getDirectionFile(); err != nil {
		return
	}

	directionFileSlice.Append(newUUID.String())

	if directionFileJSON, err = jsonIter.Marshal(directionFileSlice.ToSlice()); err != nil {
		return
	}

	if err = directionFile.Write(directionFileJSON, filesystems.Flag(os.O_CREATE|os.O_WRONLY), filesystems.Mode(0755)).GetError(); err != nil {
		return
	}

	messageFile = filesystems.NewFile(filesystems.Abs(directionFile.GetBasePath(), newUUID.String()))

	if messageFileContent, err = jsonIter.Marshal(map[string]string{"id": newUUID.String(), "content": form.Content}); err != nil {
		return
	}

	if err = messageFile.Write(messageFileContent, filesystems.Flag(os.O_CREATE|os.O_WRONLY), filesystems.Mode(0755)).GetError(); err != nil {
		return
	}

	return
}

// Destroy 留言板服务：删除信息
func (*MessageBoardService) Destroy(form *request.MessageBoardDestroyRequest) (err error) {
	var (
		directionFile      filesystems.Filesystem
		directionFileJSON  []byte
		directionFileSlice anySlices.AnySlicer[string]
		idx                int
	)

	if directionFile, directionFileSlice, err = New.MessageBoard().getDirectionFile(); err != nil {
		return
	}

	if idx = directionFileSlice.GetIndexByValue(form.ID); idx == -1 {
		return errors.New("ID不存在")
	}

	directionFileSlice.RemoveByIndex(idx)

	if directionFileJSON, err = jsonIter.Marshal(directionFileSlice.ToSlice()); err != nil {
		return
	}

	if err = directionFile.Write(directionFileJSON, filesystems.Flag(os.O_WRONLY|os.O_TRUNC|os.O_CREATE), filesystems.Mode(0755)).GetError(); err != nil {
		return
	}

	messageFile := filesystems.NewFile(filesystems.Abs(directionFile.GetBasePath(), form.ID))
	if err = messageFile.Remove().GetError(); err != nil {
		return
	}

	return
}
