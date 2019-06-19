package book

import (
	"fmt"
	"reading-club-backend/database/entity"
	"time"

	"github.com/jinzhu/gorm"

	"reading-club-backend/constant"
	"reading-club-backend/database"
	"reading-club-backend/dto"
)

var db *gorm.DB
var err error

// GetBookByID get book detail by unique id
func GetBookByID(bookID int) (book entity.Book, bookError *dto.BookErrorResponse) {
	var bookRsp entity.Book

	db, err = database.GetDBConnection()
	if err != nil {
		return bookRsp, &dto.BookErrorResponse{Error: constant.CanNotConnectDatabase, ErrorCode: constant.CanNotConnectDatabaseCode}
	}

	defer db.Close()

	errors := db.Where("id = ?", bookID).First(&bookRsp).GetErrors()

	for _, err := range errors {
		if gorm.IsRecordNotFoundError(err) {
			return bookRsp, &dto.BookErrorResponse{Error: constant.BookNotFound, ErrorCode: constant.BookNotFoundCode}
		}
		return bookRsp, &dto.BookErrorResponse{Error: constant.InternalServerError, ErrorCode: constant.InternalServerErrorCode}
	}

	if bookRsp.ID <= 0 {
		return bookRsp, &dto.BookErrorResponse{Error: constant.BookNotFound, ErrorCode: constant.BookNotFoundCode}
	}

	return bookRsp, nil
}

// GetBookByName get book detail by unique id
func GetBookByName(bookName string) (entity.Book, *dto.BookErrorResponse) {

	var bookRsp entity.Book
	var errorResponse dto.BookErrorResponse

	db, err = database.GetDBConnection()
	if err != nil {
		return bookRsp, &dto.BookErrorResponse{Error: constant.CanNotConnectDatabase, ErrorCode: constant.CanNotConnectDatabaseCode}
	}

	defer db.Close()

	visible := 1
	errors := db.Where("book_name = ? and visibility = ?", bookName, visible).First(&bookRsp).GetErrors()

	fmt.Println(errorResponse)
	for _, err := range errors {
		if gorm.IsRecordNotFoundError(err) {
			return bookRsp, &dto.BookErrorResponse{Error: constant.BookNotFound, ErrorCode: constant.BookNotFoundCode}
		}
		return bookRsp, &dto.BookErrorResponse{Error: constant.InternalServerError, ErrorCode: constant.InternalServerErrorCode}
	}

	return bookRsp, nil
}

// GetAllBooks : get all books
func GetAllBooks() (bookList []entity.Book, errorResponse *dto.BookErrorResponse) {
	db, err = database.GetDBConnection()
	if err != nil {
		return bookList, &dto.BookErrorResponse{Error: constant.CanNotConnectDatabase, ErrorCode: constant.CanNotConnectDatabaseCode}
	}

	defer db.Close()
	visible := 1
	errors := db.Where("visibility = ?", visible).Find(&bookList).GetErrors()

	for _, err := range errors {
		errorResponse.Error = err.Error()
		errorResponse.ErrorCode = constant.InternalServerErrorCode
		return bookList, &dto.BookErrorResponse{Error: constant.CanNotConnectDatabase, ErrorCode: constant.CanNotConnectDatabaseCode}
	}

	return bookList, nil
}

// SaveOrUpdate save or update book entity
func SaveOrUpdate(book *entity.Book) {
	db, err = database.GetDBConnection()
	if err != nil {
		fmt.Println(err)
		return
	}
	db.Save(book)
}

// DeleteBookByID delete a book by id
func DeleteBookByID(bookID int) {
	book, err := GetBookByID(bookID)
	if err != nil {
		fmt.Println(err)
		return
	}

	book.Visibility = 0
	book.UpdatedTime = time.Now()
	SaveOrUpdate(&book)
}
