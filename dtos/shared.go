package dtos

import (
	"fmt"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v8"
)

//BaseDto ---
type BaseDto struct {
	Success      bool     `json:"success"`
	FullMessages []string `json:"full_messages"`
}

//ErrorDto --
type ErrorDto struct {
	BaseDto
	Errors map[string]interface{} `json:"errors"`
}

//CreatePageMeta ---
func CreatePageMeta(request *http.Request, loadedItemsCount, page, pageize, totalItemsCount int) map[string]interface{} {
	pageMeta := map[string]interface{}{}
	pageMeta["offset"] = (page - 1) * pageize
	pageMeta["requested_page_size"] = pageize
	pageMeta["current_page_number"] = page
	pageMeta["current_items_count"] = loadedItemsCount

	pageMeta["prev_page_number"] = 1
	totalPagesCount := int(math.Ceil(float64(totalItemsCount) / float64(pageize)))
	pageMeta["total_pages_count"] = totalPagesCount

	if page < totalPagesCount {
		pageMeta["has_next_page"] = true
		pageMeta["next_page_number"] = page + 1
	} else {
		pageMeta["has_next_page"] = false
		pageMeta["next_page_number"] = 1
	}
	if page > 1 {
		pageMeta["prev_page_number"] = page - 1
	} else {
		pageMeta["has_prev_page"] = false
		pageMeta["prev_page_number"] = 1
	}

	pageMeta["next_page_url"] = fmt.Sprintf("%v?page=%d&page_size=%d", request.URL.Path, pageMeta["next_page_number"], pageMeta["requested_page_size"])
	pageMeta["prev_page_url"] = fmt.Sprintf("%s?page=%d&page_size=%d", request.URL.Path, pageMeta["prev_page_number"], pageMeta["requested_page_size"])

	response := gin.H{
		"success":  true,
		"pageMeta": pageMeta,
	}

	return response
}

//CreatePagedResponse --
func CreatePagedResponse(request *http.Request, resources []interface{}, resourceName string, page, pageSize, totalItemsCount int) map[string]interface{} {

	response := CreatePageMeta(request, len(resources), page, pageSize, totalItemsCount)
	response[resourceName] = resources
	return response
}

//CreateDetailedErrorDto --
func CreateDetailedErrorDto(key string, err error) map[string]interface{} {
	return map[string]interface{}{
		"success":       false,
		"full_messages": []string{fmt.Sprintf("s -> %v", key, err.Error())},
		"errors":        err,
	}
}

//CreateErrorDtoWithMessage --
func CreateErrorDtoWithMessage(message string) map[string]interface{} {
	return map[string]interface{}{
		"success":       false,
		"full_messages": []string{message},
	}
}

//CreateBadRequestErrorDto --- This should only be called when we have an Error that is returned from a ShouldBind which contains a lot of information
// other kind of errors should use other functions such as CreateDetailedErrorDto
func CreateBadRequestErrorDto(err error) ErrorDto {
	res := ErrorDto{}
	res.Errors = make(map[string]interface{})
	errs := err.(validator.ValidationErrors)
	res.FullMessages = make([]string, len(errs))
	count := 0
	for _, v := range errs {
		if v.ActualTag == "required" {
			var message = fmt.Sprintf("%v is required", v.Field)
			res.Errors[v.Field] = message
			res.FullMessages[count] = message
		} else {
			var message = fmt.Sprintf("%v has to be %v", v.Field, v.ActualTag)
			res.Errors[v.Field] = message
			res.FullMessages = append(res.FullMessages, message)
		}
		count++
	}
	return res
}

//CreateSuccessDto ---
func CreateSuccessDto(result map[string]interface{}) map[string]interface{} {
	result["success"] = true
	return result
}

//CreateSuccessWithMessageDto --
func CreateSuccessWithMessageDto(message string) interface{} {
	return CreateSuccessWithMessagesDto([]string{message})
}

//CreateSuccessWithMessagesDto --
func CreateSuccessWithMessagesDto(messages []string) interface{} {
	return gin.H{
		"success":       true,
		"full_messages": messages,
	}
}

//CreateSuccessWithDtoAndMessagesDto ---
func CreateSuccessWithDtoAndMessagesDto(data map[string]interface{}, messages []string) map[string]interface{} {
	data["success"] = true
	data["full_messages"] = messages
	return data
}

//CreateSuccessWithDtoAndMessageDto ---
func CreateSuccessWithDtoAndMessageDto(data map[string]interface{}, message string) map[string]interface{} {
	return CreateSuccessWithDtoAndMessagesDto(data, []string{message})
}
