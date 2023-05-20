package queries

const Items = `
		select 
			cast(i.ID as char(36)) ID, 
			i.Name, 
			cast(i.TemplateID as char(36)) as TemplateID, 
			cast(i.ParentID as char(36)) as ParentID, 
			cast(i.MasterID as char(36)) as MasterID, 
			i.Created, 
			i.Updated
		from Items i
		order by i.Name
`

const ItemsByTemplate = `
	with ItemSelect (ID, Name, TemplateID, ParentID, MasterID, Created, Updated)
	as
	(
		select ID, Name, TemplateID, ParentID, MasterID, Created, Updated 
		from Items
				where TemplateID in (%s)
		UNION ALL
		select rec.ID, rec.Name, rec.TemplateID, rec.ParentID, rec.MasterID, rec.Created, rec.Updated
		from Items rec
			inner join ItemSelect its
			on its.ParentID = rec.ID
	)
	select distinct
		cast(i.ID as char(36)) as ID,
		i.Name,
		cast(i.TemplateID as char(36)) as TemplateID,
		cast(i.ParentID as char(36)) as ParentID,
		cast(i.MasterID as char(36)) as MasterID,
		i.Created,
		i.Updated
	
	from ItemSelect i;`

const TemplatesByRoot = `
with ChildSelect (ID, Name, TemplateID, ParentID, MasterID, Type, BaseTemplates, StandardValuesField, Shared, Unversioned)
as
(
    select i.ID, i.Name, i.TemplateID, i.ParentID, i.MasterID, cast('' as nvarchar(max)) as Type, cast('' as nvarchar(max)) as BaseTemplates, cast('' as nvarchar(max)) as StandardValuesId, cast('0' as nvarchar(max)) as Shared, cast('0' as nvarchar(max)) as Unversioned
    from Items i
    where i.ParentID = '%s'

    UNION ALL
    select rec.ID, rec.Name, rec.TemplateID, rec.ParentID, rec.MasterID,  isnull(Replace(Replace(sf.Value, '{',''), '}', ''), '') as Type, isnull(Replace(Replace(UPPER(b.Value), '{',''), '}', ''), '') as BaseTemplates, isnull(Replace(Replace(UPPER(sv.[Value]), '{',''), '}', ''), '') as StandardValuesId, isnull(sh.Value, '0') as Shared, isnull(unv.Value, '0') as Unversioned
    from Items rec
        OUTER APPLY (select sf.Value from SharedFields sf where sf.ItemId = rec.ID and sf.FieldId = 'AB162CC0-DC80-4ABF-8871-998EE5D7BA32') sf
        OUTER APPLY (select b.Value from SharedFields b where b.ItemId = rec.ID and b.FieldId = '12C33F3F-86C5-43A5-AEB4-5598CEC45116') b
        OUTER APPLY (select sv.Value from SharedFields sv where sv.ItemId = rec.ID and sv.FieldId = 'F7D48A55-2158-4F02-9356-756654404F73') sv
        OUTER APPLY (select sh.Value from SharedFields sh where sh.ItemId = rec.ID and sh.FieldId = 'BE351A73-FCB0-4213-93FA-C302D8AB4F51') sh
        OUTER APPLY (select unv.Value from SharedFields unv where unv.ItemId = rec.ID and unv.FieldId = '39847666-389D-409B-95BD-F2016F11EED5') unv
    inner join ChildSelect t
        on t.ParentID = rec.ID
), ParentSelect (ID, Name, TemplateID, ParentID, MasterID, Type, BaseTemplates, StandardValuesField, Shared, Unversioned)
as
(
    select i.ID, i.Name, i.TemplateID, i.ParentID, i.MasterID, cast('' as nvarchar(max)) as Type, cast('' as nvarchar(max)) as BaseTemplates, cast('' as nvarchar(max)) as StandardValuesId, cast('0' as nvarchar(max)) as Shared, cast('0' as nvarchar(max)) as Unversioned
    from Items i
    where i.ParentID = '%s'

    UNION ALL
    select rec.ID, rec.Name, rec.TemplateID, rec.ParentID, rec.MasterID,  isnull(Replace(Replace(sf.Value, '{',''), '}', ''), '') as Type, isnull(Replace(Replace(UPPER(b.Value), '{',''), '}', ''), '') as BaseTemplates, isnull(Replace(Replace(UPPER(sv.[Value]), '{',''), '}', ''), '') as StandardValuesId, isnull(sh.Value, '0') as Shared, isnull(unv.Value, '0') as Unversioned
    from Items rec
        OUTER APPLY (select sf.Value from SharedFields sf where sf.ItemId = rec.ID and sf.FieldId = 'AB162CC0-DC80-4ABF-8871-998EE5D7BA32') sf
        OUTER APPLY (select b.Value from SharedFields b where b.ItemId = rec.ID and b.FieldId = '12C33F3F-86C5-43A5-AEB4-5598CEC45116') b
        OUTER APPLY (select sv.Value from SharedFields sv where sv.ItemId = rec.ID and sv.FieldId = 'F7D48A55-2158-4F02-9356-756654404F73') sv
        OUTER APPLY (select sh.Value from SharedFields sh where sh.ItemId = rec.ID and sh.FieldId = 'BE351A73-FCB0-4213-93FA-C302D8AB4F51') sh
        OUTER APPLY (select unv.Value from SharedFields unv where unv.ItemId = rec.ID and unv.FieldId = '39847666-389D-409B-95BD-F2016F11EED5') unv
        inner join ParentSelect t
        on rec.ParentID = t.ID
)
select distinct
    cast(t.ID as char(36)) as ID,
    t.Name,
    cast(t.TemplateID as char(36)) as TemplateID,
    cast(t.ParentID as char(36)) as ParentID,
    cast(t.MasterID as char(36)) as MasterID,
    t.Type, t.BaseTemplates, t.StandardValuesField, t.Shared, t.Unversioned
from ChildSelect t
    union select distinct
    cast(t.ID as char(36)) as ID,
    t.Name,
    cast(t.TemplateID as char(36)) as TemplateID,
    cast(t.ParentID as char(36)) as ParentID,
    cast(t.MasterID as char(36)) as MasterID,
    t.Type, t.BaseTemplates, t.StandardValuesField, t.Shared, t.Unversioned
from ParentSelect t`

const FieldValuesByFieldAndTemplate = `
with FieldValues (ValueID, ItemID, FieldID, Value, Version, Language, Source)
as
(
	select
		ID, ItemId, FieldId, Value, 1, 'en', 'SharedFields'
	from SharedFields
	union
	select
		ID, ItemId, FieldId, Value, Version, Language, 'VersionedFields'
	from VersionedFields
	union
	select
		ID, ItemId, FieldId, Value, 1, Language, 'UnversionedFields'
	from UnversionedFields
)
select 
	cast(fv.ValueID as char(36)) as ValueID, 
	cast(fv.ItemID as char(36)) as ItemID, 
	f.Name, 
	cast(fv.FieldID as char(36)) as FieldID, 
	fv.Value, fv.Version, 
	fv.Language, 
	fv.Source
from
	FieldValues fv
		join Items f
			on fv.FieldID = f.ID
		join Items i
			on fv.ItemID = i.ID
where fv.FieldID in (%s) and i.TemplateID in (%s)
order by fv.Source, f.Name, fv.Language, fv.Version;
`

const FieldValuesByField = `with FieldValues (ValueID, ItemID, FieldID, Value, Version, Language, Source)
as
(
	select
		ID, ItemId, FieldId, Value, 1, 'en', 'SharedFields'
	from SharedFields
	union
	select
		ID, ItemId, FieldId, Value, Version, Language, 'VersionedFields'
	from VersionedFields
	union
	select
		ID, ItemId, FieldId, Value, 1, Language, 'UnversionedFields'
	from UnversionedFields
)
select 
	cast(fv.ValueID as char(36)) as ValueID, 
	cast(fv.ItemID as char(36)) as ItemID, 
	f.Name, 
	cast(fv.FieldID as char(36)) as FieldID, 
	fv.Value, fv.Version, 
	fv.Language, 
	fv.Source
from
	FieldValues fv
		join Items f
			on fv.FieldID = f.ID
where fv.FieldID in (%s)
order by fv.Source, f.Name, fv.Language, fv.Version;
`

const FieldValues = `
with FieldValues (ValueID, ItemID, FieldID, Value, Version, Language, Source)
as
(
	select
		ID, ItemId, FieldId, Value, 1, 'en', 'SharedFields'
	from SharedFields
	union
	select
		ID, ItemId, FieldId, Value, Version, Language, 'VersionedFields'
	from VersionedFields
	union
	select
		ID, ItemId, FieldId, Value, 1, Language, 'UnversionedFields'
	from UnversionedFields
)
select 
	cast(fv.ValueID as char(36)) as ValueID, 
	cast(fv.ItemID as char(36)) as ItemID, 
	f.Name, 
	cast(fv.FieldID as char(36)) as FieldID, 
	fv.Value, fv.Version, 
	fv.Language, 
	fv.Source
from
	FieldValues fv
		join Items f
			on fv.FieldID = f.ID
order by fv.Source, f.Name, fv.Language, fv.Version;
`

const FieldValuesMeta = `
with FieldValues (ValueID, ItemID, FieldID, Version, Language, Source)
as
(
	select
		ID, ItemId, FieldId, 1, 'en', 'SharedFields'
	from SharedFields
	union
	select
		ID, ItemId, FieldId, Version, Language, 'VersionedFields'
	from VersionedFields
	union
	select
		ID, ItemId, FieldId, 1, Language, 'UnversionedFields'
	from UnversionedFields
)
select 
	cast(fv.ValueID as char(36)) as ValueID, 
	cast(fv.ItemID as char(36)) as ItemID, 
	f.Name, 
	cast(fv.FieldID as char(36)) as FieldID, 
	fv.Version, 
	fv.Language, 
	fv.Source
from
	FieldValues fv
		join Items f
			on fv.FieldID = f.ID
order by fv.Source, f.Name, fv.Language, fv.Version;
`
