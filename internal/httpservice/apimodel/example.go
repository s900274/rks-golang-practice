package apimodel


type ExamplePost struct {
    Name string `binding:"required" example:"jim"`
    Old int `binding:"exists" example:"30"`
}

type ExampleGet struct {
    Old int
}
