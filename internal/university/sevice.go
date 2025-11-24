package university

import (
	"course-tracker/config"
	"course-tracker/internal/kafka"
	"fmt"
	"math"
	"strconv"

	"gorm.io/gorm"
)

type UniversityService struct {
	DB       *gorm.DB
	CFG      *config.Config
	Producer *kafka.Producer
}

func (s *UniversityService) GetUniversities() ([]University, error) {
	var universities []University
	if err := s.DB.Find(&universities).Error; err != nil {
		return nil, err
	}
	return universities, nil
}

func (s *UniversityService) GetFields() ([]Field, error) {
	var fields []Field
	if err := s.DB.Find(&fields).Error; err != nil {
		return nil, err
	}
	return fields, nil
}

func (s *UniversityService) GetSpecializations(fieldId string) ([]Specialization, error) {
	var specs []Specialization

	if fieldId != "" {
		if err := s.DB.Where("field_id = ?", fieldId).Find(&specs).Error; err != nil {
			return nil, err
		}
		return specs, nil
	}

	if err := s.DB.Find(&specs).Error; err != nil {
		return nil, err
	}
	return specs, nil
}

func (s *UniversityService) SearchUniversities(filters map[string]string, pageStr, limitStr string) ([]CourseDetail, int64, int, int, int, error) {
	var list []CourseDetail
	var total int64

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	qb := s.DB.Table("courses c").
		Select(`c.id, c.name, c.level, c.duration, c.course_link,
			    u.name AS university, f.name AS field, s.name AS specialization`).
		Joins("JOIN universities u ON u.id = c.university_id").
		Joins("JOIN fields f ON f.id = c.field_id").
		Joins("JOIN specializations s ON s.id = c.specialization_id")

	if filters["universityId"] != "" {
		qb = qb.Where("c.university_id = ?", filters["universityId"])
	}
	if filters["fieldId"] != "" {
		qb = qb.Where("c.field_id = ?", filters["fieldId"])
	}
	if filters["specializationId"] != "" {
		qb = qb.Where("c.specialization_id = ?", filters["specializationId"])
	}
	if filters["level"] != "" {
		qb = qb.Where("c.level = ?", filters["level"])
	}
	if filters["duration"] != "" {
		qb = qb.Where("c.duration = ?", filters["duration"])
	}
	if filters["q"] != "" {
		search := fmt.Sprintf("%%%s%%", filters["q"])
		qb = qb.Where("c.name ILIKE ?", search)
	}

	countQ := qb
	if err := countQ.Count(&total).Error; err != nil {
		return nil, 0, page, limit, 0, err
	}

	if err := qb.
		Order("c.name ASC").
		Offset(offset).
		Limit(limit).
		Scan(&list).Error; err != nil {
		return nil, 0, page, limit, 0, err
	}

	pages := int(math.Ceil(float64(total) / float64(limit)))
	return list, total, page, limit, pages, nil
}

func (s *UniversityService) AddCourse(c Course) error {
	if err := s.DB.Create(&c).Error; err != nil {
		return err
	}

	evt := kafka.CourseEvent{
		Name:             c.Name,
		UniversityID:     c.UniversityID,
		FieldID:          c.FieldID,
		SpecializationID: c.SpecializationID,
		Level:            c.Level,
		Duration:         c.Duration,
		Source:           "manual",
	}

	return s.Producer.PublishCourseEvent(evt)
}

func (s *UniversityService) GetSubscribersForCourse(evt kafka.CourseEvent) ([]string, error) {
	q := s.DB.Model(&Subscription{})

	if evt.UniversityID != nil {
		q = q.Where("university_id IS NULL OR university_id = ?", *evt.UniversityID)
	}

	if evt.FieldID != nil {
		q = q.Where("field_id IS NULL OR field_id = ?", *evt.FieldID)
	}

	if evt.SpecializationID != nil {
		q = q.Where("specialization_id IS NULL OR specialization_id = ?", *evt.SpecializationID)
	}

	if evt.Level != nil {
		q = q.Where("level IS NULL OR level = ?", *evt.Level)
	}

	if evt.Duration != nil {
		q = q.Where("duration IS NULL OR duration = ?", *evt.Duration)
	}

	var emails []string
	if err := q.Distinct("email").Pluck("email", &emails).Error; err != nil {
		return nil, err
	}

	return emails, nil
}

// func (s *UniversityService) UploadCourses(file multipart.File) error {
// 	courses, err := utils.ParseCourseExcel(file)
// 	if err != nil {
// 		return err
// 	}
// 	if len(courses) == 0 {
// 		return errors.New("no valid records found")
// 	}

// 	for _, c := range courses {
// 		err := s.Repo.SaveCourse(c)
// 		if err != nil {
// 			continue
// 		}

// 		evt := kafka.CourseEvent{
// 			Name:             c.Name,
// 			UniversityID:     c.UniversityID,
// 			FieldID:          c.FieldID,
// 			SpecializationID: c.SpecializationID,
// 			Level:            c.Level,
// 			Duration:         c.Duration,
// 			Source:           "excel_upload",
// 		}
// 		_ = s.Producer.PublishCourseEvent(evt)
// 	}
// 	return nil
// }
