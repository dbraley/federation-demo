// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package accounts

import (
	"fmt"
	"io"
	"strconv"
)

type AuthRule struct {
	And  []*AuthRule `json:"and"`
	Or   []*AuthRule `json:"or"`
	Not  *AuthRule   `json:"not"`
	Rule *string     `json:"rule"`
}

type ContainsFilter struct {
	Point   *PointRef   `json:"point"`
	Polygon *PolygonRef `json:"polygon"`
}

type CustomHTTP struct {
	URL                  string     `json:"url"`
	Method               HTTPMethod `json:"method"`
	Body                 *string    `json:"body"`
	Graphql              *string    `json:"graphql"`
	Mode                 *Mode      `json:"mode"`
	ForwardHeaders       []string   `json:"forwardHeaders"`
	SecretHeaders        []string   `json:"secretHeaders"`
	IntrospectionHeaders []string   `json:"introspectionHeaders"`
	SkipIntrospection    *bool      `json:"skipIntrospection"`
}

type DateTimeFilter struct {
	Eq      *string        `json:"eq"`
	In      []*string      `json:"in"`
	Le      *string        `json:"le"`
	Lt      *string        `json:"lt"`
	Ge      *string        `json:"ge"`
	Gt      *string        `json:"gt"`
	Between *DateTimeRange `json:"between"`
}

type DateTimeRange struct {
	Min string `json:"min"`
	Max string `json:"max"`
}

type FloatFilter struct {
	Eq      *float64    `json:"eq"`
	In      []*float64  `json:"in"`
	Le      *float64    `json:"le"`
	Lt      *float64    `json:"lt"`
	Ge      *float64    `json:"ge"`
	Gt      *float64    `json:"gt"`
	Between *FloatRange `json:"between"`
}

type FloatRange struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

type GenerateMutationParams struct {
	Add    *bool `json:"add"`
	Update *bool `json:"update"`
	Delete *bool `json:"delete"`
}

type GenerateQueryParams struct {
	Get       *bool `json:"get"`
	Query     *bool `json:"query"`
	Password  *bool `json:"password"`
	Aggregate *bool `json:"aggregate"`
}

type Int64Filter struct {
	Eq      *string     `json:"eq"`
	In      []*string   `json:"in"`
	Le      *string     `json:"le"`
	Lt      *string     `json:"lt"`
	Ge      *string     `json:"ge"`
	Gt      *string     `json:"gt"`
	Between *Int64Range `json:"between"`
}

type Int64Range struct {
	Min string `json:"min"`
	Max string `json:"max"`
}

type IntFilter struct {
	Eq      *int      `json:"eq"`
	In      []*int    `json:"in"`
	Le      *int      `json:"le"`
	Lt      *int      `json:"lt"`
	Ge      *int      `json:"ge"`
	Gt      *int      `json:"gt"`
	Between *IntRange `json:"between"`
}

type IntRange struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type IntersectsFilter struct {
	Polygon      *PolygonRef      `json:"polygon"`
	MultiPolygon *MultiPolygonRef `json:"multiPolygon"`
}

type MultiPolygon struct {
	Polygons []*Polygon `json:"polygons"`
}

type MultiPolygonRef struct {
	Polygons []*PolygonRef `json:"polygons"`
}

type NearFilter struct {
	Distance   float64   `json:"distance"`
	Coordinate *PointRef `json:"coordinate"`
}

type Point struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type PointGeoFilter struct {
	Near   *NearFilter   `json:"near"`
	Within *WithinFilter `json:"within"`
}

type PointList struct {
	Points []*Point `json:"points"`
}

type PointListRef struct {
	Points []*PointRef `json:"points"`
}

type PointRef struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type Polygon struct {
	Coordinates []*PointList `json:"coordinates"`
}

type PolygonGeoFilter struct {
	Near       *NearFilter       `json:"near"`
	Within     *WithinFilter     `json:"within"`
	Contains   *ContainsFilter   `json:"contains"`
	Intersects *IntersectsFilter `json:"intersects"`
}

type PolygonRef struct {
	Coordinates []*PointListRef `json:"coordinates"`
}

type StringExactFilter struct {
	Eq      *string      `json:"eq"`
	In      []*string    `json:"in"`
	Le      *string      `json:"le"`
	Lt      *string      `json:"lt"`
	Ge      *string      `json:"ge"`
	Gt      *string      `json:"gt"`
	Between *StringRange `json:"between"`
}

type StringFullTextFilter struct {
	Alloftext *string `json:"alloftext"`
	Anyoftext *string `json:"anyoftext"`
}

type StringHashFilter struct {
	Eq *string   `json:"eq"`
	In []*string `json:"in"`
}

type StringRange struct {
	Min string `json:"min"`
	Max string `json:"max"`
}

type StringRegExpFilter struct {
	Regexp *string `json:"regexp"`
}

type StringTermFilter struct {
	Allofterms *string `json:"allofterms"`
	Anyofterms *string `json:"anyofterms"`
}

type User struct {
	ID       string  `json:"id"`
	Name     *string `json:"name"`
	Username *string `json:"username"`
}

func (User) IsEntity() {}

type WithinFilter struct {
	Polygon *PolygonRef `json:"polygon"`
}

type DgraphIndex string

const (
	DgraphIndexInt      DgraphIndex = "int"
	DgraphIndexInt64    DgraphIndex = "int64"
	DgraphIndexFloat    DgraphIndex = "float"
	DgraphIndexBool     DgraphIndex = "bool"
	DgraphIndexHash     DgraphIndex = "hash"
	DgraphIndexExact    DgraphIndex = "exact"
	DgraphIndexTerm     DgraphIndex = "term"
	DgraphIndexFulltext DgraphIndex = "fulltext"
	DgraphIndexTrigram  DgraphIndex = "trigram"
	DgraphIndexRegexp   DgraphIndex = "regexp"
	DgraphIndexYear     DgraphIndex = "year"
	DgraphIndexMonth    DgraphIndex = "month"
	DgraphIndexDay      DgraphIndex = "day"
	DgraphIndexHour     DgraphIndex = "hour"
	DgraphIndexGeo      DgraphIndex = "geo"
)

var AllDgraphIndex = []DgraphIndex{
	DgraphIndexInt,
	DgraphIndexInt64,
	DgraphIndexFloat,
	DgraphIndexBool,
	DgraphIndexHash,
	DgraphIndexExact,
	DgraphIndexTerm,
	DgraphIndexFulltext,
	DgraphIndexTrigram,
	DgraphIndexRegexp,
	DgraphIndexYear,
	DgraphIndexMonth,
	DgraphIndexDay,
	DgraphIndexHour,
	DgraphIndexGeo,
}

func (e DgraphIndex) IsValid() bool {
	switch e {
	case DgraphIndexInt, DgraphIndexInt64, DgraphIndexFloat, DgraphIndexBool, DgraphIndexHash, DgraphIndexExact, DgraphIndexTerm, DgraphIndexFulltext, DgraphIndexTrigram, DgraphIndexRegexp, DgraphIndexYear, DgraphIndexMonth, DgraphIndexDay, DgraphIndexHour, DgraphIndexGeo:
		return true
	}
	return false
}

func (e DgraphIndex) String() string {
	return string(e)
}

func (e *DgraphIndex) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = DgraphIndex(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid DgraphIndex", str)
	}
	return nil
}

func (e DgraphIndex) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type HTTPMethod string

const (
	HTTPMethodGet    HTTPMethod = "GET"
	HTTPMethodPost   HTTPMethod = "POST"
	HTTPMethodPut    HTTPMethod = "PUT"
	HTTPMethodPatch  HTTPMethod = "PATCH"
	HTTPMethodDelete HTTPMethod = "DELETE"
)

var AllHTTPMethod = []HTTPMethod{
	HTTPMethodGet,
	HTTPMethodPost,
	HTTPMethodPut,
	HTTPMethodPatch,
	HTTPMethodDelete,
}

func (e HTTPMethod) IsValid() bool {
	switch e {
	case HTTPMethodGet, HTTPMethodPost, HTTPMethodPut, HTTPMethodPatch, HTTPMethodDelete:
		return true
	}
	return false
}

func (e HTTPMethod) String() string {
	return string(e)
}

func (e *HTTPMethod) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = HTTPMethod(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid HTTPMethod", str)
	}
	return nil
}

func (e HTTPMethod) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type Mode string

const (
	ModeBatch  Mode = "BATCH"
	ModeSingle Mode = "SINGLE"
)

var AllMode = []Mode{
	ModeBatch,
	ModeSingle,
}

func (e Mode) IsValid() bool {
	switch e {
	case ModeBatch, ModeSingle:
		return true
	}
	return false
}

func (e Mode) String() string {
	return string(e)
}

func (e *Mode) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Mode(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Mode", str)
	}
	return nil
}

func (e Mode) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
