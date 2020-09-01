package model

import (
	"fmt"
	"strconv"

	"actiontech.cloud/universe/sqle/v3/sqle/errors"
	"github.com/jinzhu/gorm"
)

const (
	RULE_LEVEL_NORMAL = "normal"
	RULE_LEVEL_NOTICE = "notice"
	RULE_LEVEL_WARN   = "warn"
	RULE_LEVEL_ERROR  = "error"
)

var RuleLevelMap = map[string]int{
	RULE_LEVEL_NORMAL: 0,
	RULE_LEVEL_NOTICE: 1,
	RULE_LEVEL_WARN:   2,
	RULE_LEVEL_ERROR:  3,
}

type RuleTemplate struct {
	Model
	Name      string     `json:"name"`
	Desc      string     `json:"desc"`
	Rules     []Rule     `json:"-" gorm:"many2many:rule_template_rule"`
	Instances []Instance `json:"instance_list" gorm:"many2many:instance_rule_template"`
}

type Rule struct {
	Name  string `json:"name" gorm:"primary_key"`
	Desc  string `json:"desc"`
	Value string `json:"value"`
	Level string `json:"level" example:"error"` // notice, warn, error
}

func (r Rule) TableName() string {
	return "rules"
}

func (r *Rule) GetValue() string {
	if r == nil {
		return ""
	}
	return r.Value
}

func (r *Rule) GetValueInt(defaultRule *Rule) int64 {
	value := r.GetValue()
	i, err := strconv.ParseInt(value, 10, 64)
	if err == nil {
		return i
	}
	i, err = strconv.ParseInt(defaultRule.GetValue(), 10, 64)
	if err == nil {
		return i
	}
	return 0
}

// RuleTemplateDetail use for http request and swagger docs;
// it is same as RuleTemplate, but display Rules in json format.
type RuleTemplateDetail struct {
	RuleTemplate
	Rules     []Rule     `json:"rule_list"`
	Instances []Instance `json:"instance_list"`
}

func (r *RuleTemplate) Detail() RuleTemplateDetail {
	data := RuleTemplateDetail{
		RuleTemplate: *r,
		Rules:        r.Rules,
		Instances:    r.Instances,
	}
	if r.Rules == nil {
		data.Rules = []Rule{}
	}
	if r.Instances == nil {
		data.Instances = []Instance{}
	}
	return data
}

func (s *Storage) GetTemplateById(templateId string) (RuleTemplate, bool, error) {
	t := RuleTemplate{}
	err := s.db.Preload("Rules").Preload("Instances").Where("id = ?", templateId).First(&t).Error
	if err == gorm.ErrRecordNotFound {
		return t, false, nil
	}
	return t, true, errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) GetTemplateByName(name string) (RuleTemplate, bool, error) {
	t := RuleTemplate{}
	err := s.db.Preload("Rules").Where("name = ?", name).First(&t).Error
	if err == gorm.ErrRecordNotFound {
		return t, false, nil
	}
	return t, true, errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) UpdateTemplateRules(tpl *RuleTemplate, rules ...Rule) error {
	err := s.db.Model(tpl).Association("Rules").Replace(rules).Error
	return errors.New(errors.CONNECT_STORAGE_ERROR, err)
}
func (s *Storage) UpdateTemplateInst(tpl *RuleTemplate, inst ...Instance) error {
	err := s.db.Model(tpl).Association("Instances").Replace(inst).Error
	return errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) GetAllTemplate() ([]RuleTemplate, error) {
	ts := []RuleTemplate{}
	err := s.db.Preload("Instances").Find(&ts).Error
	return ts, errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func GetRuleMapFromAllArray(allRules ...[]Rule) map[string]Rule {
	ruleMap := map[string]Rule{}
	for _, rules := range allRules {
		for _, rule := range rules {
			ruleMap[rule.Name] = rule
		}
	}
	return ruleMap
}

func (s *Storage) GetAllRule() ([]Rule, error) {
	rules := []Rule{}
	err := s.db.Find(&rules).Error
	return rules, errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) GetRulesByInstanceId(instanceId string) ([]Rule, error) {
	rules := []Rule{}
	instance, _, err := s.GetInstById(instanceId)
	if err != nil {
		return rules, errors.New(errors.CONNECT_STORAGE_ERROR, err)
	}
	templates := instance.RuleTemplates
	if len(templates) <= 0 {
		return rules, nil
	}
	templateIds := make([]string, len(templates))
	for n, template := range templates {
		templateIds[n] = fmt.Sprintf("%v", template.ID)
	}

	err = s.db.Table("rules").Select("rules.name, rules.value, rules.level").
		Joins("inner join rule_template_rule on rule_template_rule.rule_name = rules.name").
		Where("rule_template_rule.rule_template_id in (?)", templateIds).Scan(&rules).Error
	return rules, errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) UpdateRuleValueByName(name, value string) error {
	err := s.db.Table("rules").Where("name = ?", name).
		Update(map[string]string{"value": value}).Error
	return errors.New(errors.CONNECT_STORAGE_ERROR, err)
}

func (s *Storage) GetInstancesNameByTemplate(tpl RuleTemplate) ([]string, error) {
	names := []string{}
	rows, err := s.db.Table("instances").Select("instances.name").
		Joins("inner join instance_rule_template on instance_rule_template.instance_id = instances.id").
		Where("instances.deleted_at IS NULL and instance_rule_template.rule_template_id = ?", tpl.ID).Rows()
	if err != nil {
		return nil, errors.New(errors.CONNECT_STORAGE_ERROR, err)
	}
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return nil, errors.New(errors.CONNECT_STORAGE_ERROR, err)
		}
		names = append(names, name)
	}
	return names, nil
}