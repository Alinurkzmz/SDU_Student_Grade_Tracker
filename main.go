package main

import (
  "fmt"
  "strconv"

  "github.com/go-redis/redis"
)

type Student struct {
  ID     string
  Grades map[string]int
}

type Application struct {
  redisClient *redis.Client
}

func NewApplication() *Application {
  client := redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
  })
  return &Application{
    redisClient: client,
  }
}

func (app *Application) AddGrade(studentID, courseCode string, grade int) error {
  student, err := app.getStudent(studentID)
  if err != nil {
    return err
  }

  student.Grades[courseCode] = grade
  return app.saveStudent(student)
}

func (app *Application) ViewGrades(studentID string) (map[string]int, error) {
  student, err := app.getStudent(studentID)
  if err != nil {
    return nil, err
  }

  return student.Grades, nil
}

func (app *Application) UpdateGrade(studentID, courseCode string, newGrade int) error {
  student, err := app.getStudent(studentID)
  if err != nil {
    return err
  }

  student.Grades[courseCode] = newGrade
  return app.saveStudent(student)
}

func (app *Application) DeleteGrade(studentID, courseCode string) error {
  studentKey := fmt.Sprintf("student:%s", studentID)
  gradeKey := fmt.Sprintf("grade:%s", courseCode)

  err := app.redisClient.HDel(studentKey, gradeKey).Err()
  if err != nil {
    return err
  }

  err = app.redisClient.HDel(gradeKey, studentKey).Err()
  if err != nil {
    return err
  }

  return nil
}

func (app *Application) getStudent(studentID string) (*Student, error) {
  studentKey := fmt.Sprintf("student:%s", studentID)

  grades, err := app.redisClient.HGetAll(studentKey).Result()
  if err != nil {
    return nil, err
  }

  gradesMap := make(map[string]int)
  for courseCode, gradeStr := range grades {
    grade, _ := strconv.Atoi(gradeStr)
    gradesMap[courseCode] = grade
  }

  return &Student{
    ID:     studentID,
    Grades: gradesMap,
  }, nil
}

func (app *Application) saveStudent(student *Student) error {
  studentKey := fmt.Sprintf("student:%s", student.ID)

  grades := make(map[string]interface{})
  for courseCode, grade := range student.Grades {
    grades[courseCode] = grade
  }

  err := app.redisClient.HMSet(studentKey, grades).Err()
  if err != nil {
    return err
  }

  return nil
}

func main() {
  app := NewApplication()

  studentID := "2456"
  courseCode := "CSC103"
  grade := 90

  err := app.AddGrade(studentID, courseCode, grade)
  if err != nil {
    fmt.Println("Error adding grade:", err)
    return
  }

  grades, err := app.ViewGrades(studentID)
  if err != nil {
    fmt.Println("Error viewing grades:", err)
    return
  }

  fmt.Println("Grades:", grades)

  err = app.UpdateGrade(studentID, courseCode, 96)
  if err != nil {
    fmt.Println("Error updating grade:", err)
    return
  }

  err = app.DeleteGrade(studentID, courseCode)
  if err != nil {
    fmt.Println("Error deleting grade:", err)
    return
  }
}