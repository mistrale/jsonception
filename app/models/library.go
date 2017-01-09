package models

// Library : test container
type Library struct {
	LibraryID int    `json:"library_id" gorm:"primary_key"`
	Name      string `json:"name"`
	TestIDs   []int  `json:"test_ids" sql:"-"`
	Tests     []Test `json:"tests" db:"-" gorm:"many2many:library_tests;"`
	Uuid      string `json:"-" db:"-"`
}

//
// type IntArray []int
//
// // Value returns the driver compatible value
// func (a IntArray) Value() (driver.Value, error) {
// 	var strs []string
// 	for _, i := range a {
// 		strs = append(strs, strconv.Itoa(i))
// 	}
// 	return "[" + strings.Join(strs, ",") + "]", nil
// }
//
// func strToIntSlice(s string) IntArray {
// 	r := strings.Trim(s, "[]")
// 	a := make(IntArray, 0, 10)
// 	for _, t := range strings.Split(r, ",") {
// 		i, _ := strconv.Atoi(t)
// 		a = append(a, i)
// 	}
// 	return a
// }
//
// func (s *IntArray) Scan(src interface{}) error {
// 	fmt.Printf("ON EST ICI")
// 	asBytes, ok := src.([]byte)
// 	if !ok {
// 		return error(errors.New("Scan source was not []byte"))
// 	}
// 	asString := string(asBytes)
// 	(*s) = strToIntSlice(asString)
// 	return nil
// }
