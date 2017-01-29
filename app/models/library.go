package models

// Library : test container
type Library struct {
	LibraryID int    `json:"library_id" gorm:"primary_key"`
	Name      string `json:"name"`
	TestIDs   []int  `json:"test_ids" sql:"-"`
	Tests     []Test `json:"tests" gorm:"many2many:library_tests;"`
	Uuid      string `json:"-" db:"-"`
}
