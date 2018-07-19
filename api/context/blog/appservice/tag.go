package appservice

import (
	account "lmm/api/context/account/domain/model"
	"lmm/api/context/blog/domain"
	"lmm/api/context/blog/domain/factory"
	"lmm/api/utils/strings"
)

func (app *AppService) AddNewTagToBlog(user *account.User, blogIDStr, tagName string) error {
	blogID, err := strings.StrToUint64(blogIDStr)
	if err != nil {
		return domain.ErrNoSuchBlog
	}
	tag, err := factory.NewTag(blogID, tagName)
	if err != nil {
		return err
	}
	if err := app.tagRepository.Add(tag); err != nil {
		return err
	}
	return nil
}
