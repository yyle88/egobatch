package example2_test

import (
	"context"
	"math/rand/v2"
	"testing"

	"github.com/yyle88/egobatch"
	"github.com/yyle88/egobatch/erxgroup"
	"github.com/yyle88/egobatch/internal/examples/example2"
	"github.com/yyle88/egobatch/internal/myerrors"
	"github.com/yyle88/must"
	"github.com/yyle88/neatjson/neatjsons"
)

func TestRun(t *testing.T) {
	ctx := context.Background()
	classes := example2.NewClasses(5)
	taskResults := processClasses(ctx, classes)
	classesStudentsScores := taskResults.Flatten(func(arg *example2.Class, err *myerrors.Error) *example2.ClassStudentsScores {
		return &example2.ClassStudentsScores{
			Class:          arg,
			StudentsScores: nil,
			Erx:            err,
		}
	})
	t.Log(neatjsons.S(classesStudentsScores))
}

func processClasses(ctx context.Context, classes []*example2.Class) egobatch.Tasks[*example2.Class, *example2.ClassStudentsScores, *myerrors.Error] {
	taskBatch := egobatch.NewTaskBatch[*example2.Class, *example2.ClassStudentsScores, *myerrors.Error](classes)
	taskBatch.SetGlide(true)
	taskBatch.SetWaCtx(func(err error) *myerrors.Error {
		return myerrors.ErrorWrongContext("wrong-ctx-can-not-invoke-process-class-func. error=%v", err)
	})
	ego := erxgroup.NewGroup[*myerrors.Error](ctx)
	ego.SetLimit(3)
	taskBatch.EgoRun(ego, processClassFunc)
	must.Null(ego.Wait())
	return taskBatch.Tasks
}

func processClassFunc(ctx context.Context, arg *example2.Class) (*example2.ClassStudentsScores, *myerrors.Error) {
	if rand.IntN(100) < 30 {
		return nil, myerrors.ErrorServiceError("wrong-db")
	}
	studentCount := 1 + rand.IntN(5)
	students := example2.NewStudents(studentCount)
	taskResults := processStudents(ctx, students)
	studentsScores := taskResults.Flatten(func(arg *example2.Student, err *myerrors.Error) *example2.StudentScores {
		return &example2.StudentScores{
			Student:  arg,
			Scores:   nil,
			AvgScore: 0,
			Erx:      err,
		}
	})

	okCnt := 0
	okSum := float64(0)
	for _, studentScores := range studentsScores {
		if studentScores.Erx != nil {
			continue
		}
		okCnt++
		okSum += studentScores.AvgScore
	}
	avgScore := float64(0)
	if okCnt > 0 {
		avgScore = okSum / float64(okCnt)
	}

	return &example2.ClassStudentsScores{
		Class:          arg,
		StudentsScores: studentsScores,
		AvgScore:       avgScore,
		Erx:            nil,
	}, nil
}

func processStudents(ctx context.Context, students []*example2.Student) egobatch.Tasks[*example2.Student, *example2.StudentScores, *myerrors.Error] {
	taskBatch := egobatch.NewTaskBatch[*example2.Student, *example2.StudentScores, *myerrors.Error](students)
	taskBatch.SetGlide(true)
	taskBatch.SetWaCtx(func(err error) *myerrors.Error {
		return myerrors.ErrorWrongContext("wrong-ctx-can-not-invoke-process-student-func. error=%v", err)
	})
	ego := erxgroup.NewGroup[*myerrors.Error](ctx)
	ego.SetLimit(2)
	taskBatch.EgoRun(ego, processStudentFunc)
	must.Null(ego.Wait())
	return taskBatch.Tasks
}

func processStudentFunc(ctx context.Context, arg *example2.Student) (*example2.StudentScores, *myerrors.Error) {
	if rand.IntN(100) < 30 {
		return nil, myerrors.ErrorServiceError("wrong-db")
	}
	subjectCount := 1 + rand.IntN(2)
	subjects := example2.NewSubjects(subjectCount)

	taskResults := processSubjects(ctx, subjects)
	scores := taskResults.Flatten(func(arg *example2.Subject, err *myerrors.Error) *example2.SubjectScore {
		return &example2.SubjectScore{
			Subject: arg,
			Score:   0,
			Erx:     err,
		}
	})

	okCnt := 0
	okSum := float64(0)
	for _, score := range scores {
		if score.Erx != nil {
			continue
		}
		okCnt++
		okSum += float64(score.Score)
	}
	avgScore := float64(0)
	if okCnt > 0 {
		avgScore = okSum / float64(okCnt)
	}

	return &example2.StudentScores{
		Student:  arg,
		Scores:   scores,
		AvgScore: avgScore,
		Erx:      nil,
	}, nil
}

func processSubjects(ctx context.Context, subjects []*example2.Subject) egobatch.Tasks[*example2.Subject, *example2.SubjectScore, *myerrors.Error] {
	taskBatch := egobatch.NewTaskBatch[*example2.Subject, *example2.SubjectScore, *myerrors.Error](subjects)
	taskBatch.SetGlide(true)
	taskBatch.SetWaCtx(func(err error) *myerrors.Error {
		return myerrors.ErrorWrongContext("wrong-ctx-can-not-invoke-process-subject-func. error=%v", err)
	})
	ego := erxgroup.NewGroup[*myerrors.Error](ctx)
	ego.SetLimit(2)
	taskBatch.EgoRun(ego, processSubjectFunc)
	must.Null(ego.Wait())
	return taskBatch.Tasks
}

func processSubjectFunc(ctx context.Context, arg *example2.Subject) (*example2.SubjectScore, *myerrors.Error) {
	if rand.IntN(100) < 30 {
		return nil, myerrors.ErrorServiceError("wrong-db")
	}
	return &example2.SubjectScore{
		Subject: arg,
		Score:   rand.IntN(100),
		Erx:     nil,
	}, nil
}
