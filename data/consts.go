package data

import (
	"github.com/google/uuid"
)

// items
var RootID = uuid.Must(uuid.Parse("00000000-0000-0000-0000-000000000000"))
var EmptyID = RootID

var TemplatesRootID = uuid.Must(uuid.Parse("3c1715fe-6a13-4fcf-845f-de308ba9741d"))
var TemplateID = uuid.Must(uuid.Parse("AB86861A-6030-46C5-B394-E8F99E8B87DB"))
var TemplateFieldID = uuid.Must(uuid.Parse("455A3E98-A627-4B40-8035-E683A0331AC7"))
var TemplateSectionID = uuid.Must(uuid.Parse("E269FBB5-3750-427A-9149-7AA950B49301"))
var StandardTemplateID = uuid.Must(uuid.Parse("1930BBEB-7805-471A-A3BE-4858AC7CF696"))
var LayoutTemplateId = uuid.Must(uuid.Parse("3A45A723-64EE-4919-9D41-02FD40FD1466"))

var ControllerRenderingId = uuid.Must(uuid.Parse("2A3E91A0-7987-44B5-AB34-35C2D9DE83B9"))
var ItemRenderingId = uuid.Must(uuid.Parse("86776923-ECA5-4310-8DC0-AE65FE88D078"))
var MethodRenderingId = uuid.Must(uuid.Parse("39587D7D-F06D-4CB4-A25E-AA7D847EDDD0"))
var SublayoutRenderingId = uuid.Must(uuid.Parse("0A98E368-CDB9-4E1E-927C-8E0C24A003FB"))
var UrlRenderingId = uuid.Must(uuid.Parse("83E993C5-C0FC-4472-86A9-2F6CFED694E4"))
var ViewRenderingId = uuid.Must(uuid.Parse("99F8905D-4A87-4EB8-9F8B-A9BEBFB3ADD6"))
var WebControlRenderingId = uuid.Must(uuid.Parse("1DDE3F02-0BD7-4779-867A-DC578ADF91EA"))
var XmlControlRenderingId = uuid.Must(uuid.Parse("B658CE99-894A-4CB1-936B-F23F17C63B5B"))
var XslRenderingId = uuid.Must(uuid.Parse("F1F1D639-4F54-40C2-8BE0-81266B392CEB"))

var RenderingParametersID = uuid.Must(uuid.Parse("8CA06D6A-B353-44E8-BC31-B528C7306971"))

// fields

var RenderingsFieldId = uuid.Must(uuid.Parse("f1a1fe9e-a60c-4ddb-a3a0-bb5b29fe732e"))
var FinalRenderingsFieldId = uuid.Must(uuid.Parse("04bf00db-f5fb-41f7-8ab7-22408372a981"))
var StandardValuesFieldId = uuid.Must(uuid.Parse("F7D48A55-2158-4F02-9356-756654404F73"))
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

// media fields

var BlobFieldId = uuid.Must(uuid.Parse("40e50ed9-ba07-4702-992e-a912738d32dc"))
var AltFieldId = uuid.Must(uuid.Parse("65885c44-8fcd-4a7f-94f1-ee63703fe193"))
var ExtensionFieldId = uuid.Must(uuid.Parse("c06867fe-9a43-4c7d-b739-48780492d06f"))
var MimeTypeFieldId = uuid.Must(uuid.Parse("6f47a0a5-9c94-4b48-abeb-42d38def6054"))
var DescriptionFieldId = uuid.Must(uuid.Parse("ba8341a1-ff30-47b8-ae6a-f4947e4113f0"))
var HeightFieldId = uuid.Must(uuid.Parse("de2ca9e4-c117-4c8a-a139-1ff4b199d15a"))
var WidthFieldId = uuid.Must(uuid.Parse("22eac599-f13b-4607-a89d-c091763a467d"))
var CopyrightFieldId = uuid.Must(uuid.Parse("7993837d-2c79-44bc-b5ec-ae93446a68c5"))
var KeywordsFieldId = uuid.Must(uuid.Parse("2fafe7cb-2691-4800-8848-255efa1d31aa"))
var TitleFieldId = uuid.Must(uuid.Parse("3f4b20e9-36e6-4d45-a423-c86567373f82"))

var VersionedBlobFieldId = uuid.Must(uuid.Parse("dbbe7d99-1388-4357-bb34-ad71edf18ed3"))
var VersionedAltFieldId = uuid.Must(uuid.Parse("8cf45d0e-add4-4772-911a-ac6fc50f9c7d"))
var VersionedExtensionFieldId = uuid.Must(uuid.Parse("3eb149f8-de14-4220-a8f4-9e723cfae5d9"))
var VersionedMimeTypeFieldId = uuid.Must(uuid.Parse("aba5ce58-6d8f-43a4-b4fa-b181651347dd"))
var VersionedDescriptionFieldId = uuid.Must(uuid.Parse("ebeb197c-376e-47c4-95d7-7fc26682d12e"))
var VersionedHeightFieldId = uuid.Must(uuid.Parse("611ba5dd-8d26-4c9b-a95d-900ac94cb32a"))
var VersionedWidthFieldId = uuid.Must(uuid.Parse("d8e619c3-2d16-47f4-ad3b-c68103df86c2"))
var VersionedCopyrightFieldId = uuid.Must(uuid.Parse("ed74c2c9-20b3-47a4-b021-a652e1e808aa"))
var VersionedKeywordsFieldId = uuid.Must(uuid.Parse("1cc29119-82e2-449f-b208-33790d434b95"))
var VersionedTitleFieldId = uuid.Must(uuid.Parse("e625e7da-f988-4442-a598-d21040ec9815"))
