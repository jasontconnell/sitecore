package data

import (
    "github.com/google/uuid"
)


type TemplateMeta struct {
    Type            string
    BaseTemplateIds []uuid.UUID
}

type Template struct {
    TemplateMeta
    Item
    Fields        []TemplateField
    BaseTemplates []*Template
}

type TemplateField struct {
    Item ItemNode
    Type string
}