# gormery

### Description
small library for simplify building complex SQL queries in [gorm.io](https://gorm.io).

### Install
```
go get github.com/nightlord189/gormery
```

### Use
```Go
import "github.com/nightlord189/gormery"

result := make([]model.Entity, 0)

queryElems := make([]gormery.ConditionElement, 0)
queryElems = append(queryElems, gormery.Equal("id", 18))
queryElems = append(queryElems, gormery.NotEqual("parent_id", "1834"))
queryElems = append(queryElems, gormery.Like("name", "%orange%"))
queryElems = append(queryElems, gormery.MoreOrEqual("amount", 201500))

sql, elems := gormery.CombineSimpleQuery(queryElems, "AND")

d.DB.Where(sql, elems...).Find(&result).Error
```

this query will be translated to:
```SQL
SELECT * FROM entities 
WHERE 
id = 18 
AND parent_id <> '1834' 
AND name LIKE '%orange%' 
AND amount >= 201500
```

### Future features
+ combining AND, OR and other
+ brackets