// Package example2 demonstrates nested batch processing patterns
// Shows class, student, and subject score aggregation with error handling
//
// 包 example2 演示嵌套批量处理模式
// 展示班级、学生和科目成绩聚合以及错误处理
package example2

import (
	"fmt"

	"github.com/yyle88/egobatch/internal/myerrors"
)

// Class represents a class entity
// 班级实体
type Class struct {
	Name string // Class name // 班级名称
}

// ClassStudentsScores aggregates class with student scores and average
// 聚合班级及其学生成绩和平均分
type ClassStudentsScores struct {
	Class          *Class           // Class reference // 班级引用
	StudentsScores []*StudentScores // Collection of student scores // 学生成绩集合
	AvgScore       float64          // Class average score // 班级平均分
	Erx            *myerrors.Error  // Processing error if any // 处理错误（如有）
}

// Student represents a student entity
// 学生实体
type Student struct {
	Name string // Student name // 学生名称
}

// StudentScores aggregates student with subject scores and average
// 聚合学生及其科目成绩和平均分
type StudentScores struct {
	Student  *Student        // Student reference // 学生引用
	Scores   []*SubjectScore // Collection of subject scores // 科目成绩集合
	AvgScore float64         // Student average score // 学生平均分
	Erx      *myerrors.Error // Processing error if any // 处理错误（如有）
}

// Subject represents a subject entity
// 科目实体
type Subject struct {
	Name string // Subject name // 科目名称
}

// SubjectScore represents score on a subject with error tracking
// 科目成绩，包含错误跟踪
type SubjectScore struct {
	Subject *Subject        // Subject reference // 科目引用
	Score   int             // Score value // 成绩值
	Erx     *myerrors.Error // Processing error if any // 处理错误（如有）
}

// NewClasses creates a collection of class instances
// Names formatted with sequential index
//
// NewClasses 创建班级实例集合
// 名称按索引格式化
func NewClasses(classCount int) []*Class {
	classes := make([]*Class, 0, classCount)
	for idx := 0; idx < classCount; idx++ {
		classes = append(classes, &Class{
			Name: fmt.Sprintf("class(%d)", idx),
		})
	}
	return classes
}

// NewStudents creates a collection of student instances
// Names formatted with sequential index
//
// NewStudents 创建学生实例集合
// 名称按索引格式化
func NewStudents(studentCount int) []*Student {
	var students = make([]*Student, 0, studentCount)
	for idx := 0; idx < studentCount; idx++ {
		students = append(students, &Student{
			Name: fmt.Sprintf("student(%d)", idx),
		})
	}
	return students
}

// NewSubjects creates a collection of subject instances
// Names formatted with sequential index
//
// NewSubjects 创建科目实例集合
// 名称按索引格式化
func NewSubjects(subjectCount int) []*Subject {
	var subjects = make([]*Subject, 0, subjectCount)
	for idx := 0; idx < subjectCount; idx++ {
		subjects = append(subjects, &Subject{
			Name: fmt.Sprintf("subject(%d)", idx),
		})
	}
	return subjects
}
