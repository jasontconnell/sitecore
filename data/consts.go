package data

import (
	"github.com/google/uuid"
)

var RenderingsFieldId = uuid.Must(uuid.Parse("F1A1FE9E-A60C-4DDB-A3A0-BB5B29FE732E"))
var FinalRenderingsFieldId = uuid.Must(uuid.Parse("04BF00DB-F5FB-41F7-8AB7-22408372A981"))

var ControllerRenderingId = uuid.Must(uuid.Parse("2A3E91A0-7987-44B5-AB34-35C2D9DE83B9"))
var ItemRenderingId = uuid.Must(uuid.Parse("86776923-ECA5-4310-8DC0-AE65FE88D078"))
var MethodRenderingId = uuid.Must(uuid.Parse("39587D7D-F06D-4CB4-A25E-AA7D847EDDD0"))
var SublayoutRenderingId = uuid.Must(uuid.Parse("0A98E368-CDB9-4E1E-927C-8E0C24A003FB"))
var UrlRenderingId = uuid.Must(uuid.Parse("83E993C5-C0FC-4472-86A9-2F6CFED694E4"))
var ViewRenderingId = uuid.Must(uuid.Parse("99F8905D-4A87-4EB8-9F8B-A9BEBFB3ADD6"))
var WebControlRenderingId = uuid.Must(uuid.Parse("1DDE3F02-0BD7-4779-867A-DC578ADF91EA"))
var XmlControlRenderingId = uuid.Must(uuid.Parse("B658CE99-894A-4CB1-936B-F23F17C63B5B"))
var XslRenderingId = uuid.Must(uuid.Parse("F1F1D639-4F54-40C2-8BE0-81266B392CEB"))

var LayoutTemplateId = uuid.Must(uuid.Parse("3A45A723-64EE-4919-9D41-02FD40FD1466"))

var TemplateID = uuid.Must(uuid.Parse("AB86861A-6030-46C5-B394-E8F99E8B87DB"))
var TemplateFieldID = uuid.Must(uuid.Parse("455A3E98-A627-4B40-8035-E683A0331AC7"))
var TemplateSectionID = uuid.Must(uuid.Parse("E269FBB5-3750-427A-9149-7AA950B49301"))

var RenderingParametersID = uuid.Must(uuid.Parse("8CA06D6A-B353-44E8-BC31-B528C7306971"))

var RootID = uuid.Must(uuid.Parse("00000000-0000-0000-0000-000000000000"))
var EmptyID = RootID
var StandardTemplateID = uuid.Must(uuid.Parse("1930BBEB-7805-471A-A3BE-4858AC7CF696"))

var StandardValuesFieldId = uuid.Must(uuid.Parse("F7D48A55-2158-4F02-9356-756654404F73"))

var BlobFieldId = uuid.Must(uuid.Parse("40E50ED9-BA07-4702-992E-A912738D32DC"))

var FieldTypeFieldId = uuid.Must(uuid.Parse("AB162CC0-DC80-4ABF-8871-998EE5D7BA32"))
var BaseTemplatesFieldId = uuid.Must(uuid.Parse("12C33F3F-86C5-43A5-AEB4-5598CEC45116"))
var UnversionedFieldId = uuid.Must(uuid.Parse("39847666-389d-409b-95bd-f2016f11eed5"))
var SharedFieldId = uuid.Must(uuid.Parse("be351a73-fcb0-4213-93fa-c302d8ab4f51"))

var SublayoutRenderingPathFieldId = uuid.Must(uuid.Parse("e42081b6-8a95-4a11-89ce-df70ed502f57"))
var RenderingDatasourceLocationFieldId = uuid.Must(uuid.Parse("b5b27af1-25ef-405c-87ce-369b3a004016"))
var RenderingDatasourceTemplateFieldId = uuid.Must(uuid.Parse("1a7c85e5-dc0b-490d-9187-bb1dbcb4c72f"))

var ControllerRenderingControllerFieldId = uuid.Must(uuid.Parse("e64ad073-dfcc-4d20-8c0b-fe5aa6226cd7"))

var DisplayNameFieldId = uuid.Must(uuid.Parse("b5e02ad9-d56f-4c41-a065-a133db87bdeb"))

var CreatedByFieldId = uuid.Must(uuid.Parse("5dd74568-4d4b-44c1-b513-0af5f4cda34f"))
var UpdatedByFieldId = uuid.Must(uuid.Parse("badd9cf9-53e0-4d0c-bcc0-2d784c282f6a"))
var CreateDateFieldId = uuid.Must(uuid.Parse("25bed78c-4957-4165-998a-ca1b52f67497"))
var UpdateDateFieldId = uuid.Must(uuid.Parse("d9cf14b1-fa16-4ba6-9288-e8a174d4d522"))
