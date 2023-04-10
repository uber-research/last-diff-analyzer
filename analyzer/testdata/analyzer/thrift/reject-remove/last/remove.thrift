// This test file is adapted from all test files from https://github.com/thriftrw/thriftrw-go/tree/183888fb47c3c225e86d634fa3701ce8b84c1914/gen/internal/tests/thrift

struct First {}

struct Second {}

enum EnumDefault {
    Foo, Bar, Baz
}

struct PrimitiveRequiredStruct {
    1: required bool boolField
    2: required byte byteField
    3: required i16 int16Field
    4: required i32 int32Field
    5: required i64 int64Field
    6: required double doubleField
    7: required string stringField
    8: required binary binaryField
    9: required list<string> listOfStrings
    10: required set<i32> setOfInts
}

typedef map<string, string> StringMap
typedef PrimitiveRequiredStruct Primitives
typedef list<string> StringList

const list<i32> listOfInts = [1, 2, 3]

const structs.Point some_point = {"x": 1, "y": 2.0}

typedef string Key

exception InternalError {
    1: optional string message
}

service KeyValue {
    // void and no exceptions
    void setValue(1: Key key, 2: unions.ArbitraryValue value)

    void setValueV2(
        /** Key to change. */
        1: required Key key,
        /**
         * New value for the key.
         *
         * If the key already has an existing value, it will be overwritten.
         */
        2: required unions.ArbitraryValue value,
    )

    // Return with exceptions
    unions.ArbitraryValue getValue(1: Key key)
        throws (1: exceptions.DoesNotExistException doesNotExist)

    // void with exceptions
    void deleteValue(1: Key key)
        throws (
            /**
             * Raised if a value with the given key doesn't exist.
             */
            1: exceptions.DoesNotExistException doesNotExist,
            2: InternalError internalError
        )

    list<unions.ArbitraryValue> getManyValues(
        1: list<Key> range  // < reserved keyword as an argument
    ) throws (
        1: exceptions.DoesNotExistException doesNotExist,
    )

    i64 size()  // < primitve return value
}

service Cache {
    oneway void clear()
    oneway void clearAfter(1: i64 durationMS)
}

struct ConflictingNames_SetValue_Args {
    1: required string key
    2: required binary value
}

service ConflictingNames {
    void setValue(1: ConflictingNames_SetValue_Args request)
}

service non_standard_service_name {
    void non_standard_function_name()
}

typedef set<string> StringSet
typedef set<string> (go.type = "slice") StringList
typedef set<Foo> (go.type = "slice") FooList
typedef StringList MyStringList
typedef MyStringList AnotherStringList

typedef set<set<string> (go.type = "slice")> (go.type = "slice") StringListList

struct Foo {
    1: required string stringField
}

struct Bar {
    1: required set<i32> (go.type = "slice") requiredInt32ListField
    2: optional set<string> (go.type = "slice") optionalStringListField
    3: required StringList requiredTypedefStringListField
    4: optional StringList optionalTypedefStringListField
    5: required set<Foo> (go.type = "slice") requiredFooListField
    6: optional set<Foo> (go.type = "slice") optionalFooListField
    7: required FooList requiredTypedefFooListField
    8: optional FooList optionalTypedefFooListField
    9: required set<set<string> (go.type = "slice")> (go.type = "slice") requiredStringListListField
    10: required StringListList requiredTypedefStringListListField
}

const set<string> (go.type = "slice") ConstStringList = ["hello"]
const set<set<string>(go.type = "slice")> (go.type = "slice") ConstListStringList = [["hello"], ["world"]]

struct EmptyStruct {}

//////////////////////////////////////////////////////////////////////////////
// Structs with primitives

/**
 * A struct that contains primitive fields exclusively.
 *
 * All fields are required.
 */
struct PrimitiveRequiredStruct {
    1: required bool boolField
    2: required byte byteField
    3: required i16 int16Field
    4: required i32 int32Field
    5: required i64 int64Field
    6: required double doubleField
    7: required string stringField
    8: required binary binaryField
}

/**
 * A struct that contains primitive fields exclusively.
 *
 * All fields are optional.
 */
struct PrimitiveOptionalStruct {
    1: optional bool boolField
    2: optional byte byteField
    3: optional i16 int16Field
    4: optional i32 int32Field
    5: optional i64 int64Field
    6: optional double doubleField
    7: optional string stringField
    8: optional binary binaryField
}

//////////////////////////////////////////////////////////////////////////////
// Nested structs (Required)

/**
 * A point in 2D space.
 */
struct Point {
    1: required double x
    2: required double y
}

/**
 * Size of something.
 */
struct Size {
    /**
     * Width in pixels.
     */
    1: required double width
    /** Height in pixels. */
    2: required double height
}

struct Frame {
    1: required Point topLeft
    2: required Size size
}

struct Edge {
    1: required Point startPoint
    2: required Point endPoint
}

/**
 * A graph is comprised of zero or more edges.
 */
struct Graph {
    /**
     * List of edges in the graph.
     *
     * May be empty.
     */
    1: required list<Edge> edges
}

//////////////////////////////////////////////////////////////////////////////
// Nested structs (Optional)

struct ContactInfo {
    1: required string emailAddress
}

struct PersonalInfo {
    1: optional i32 age
}

struct User {
    1: required string name
    2: optional ContactInfo contact
    3: optional PersonalInfo personal
}

typedef map<string, User> UserMap

//////////////////////////////////////////////////////////////////////////////
// self-referential struct

typedef Node List

/**
 * Node is linked list of values.
 * All values are 32-bit integers.
 */
struct Node {
    1: required i32 value
    2: optional List tail
}

//////////////////////////////////////////////////////////////////////////////
// JSON tagged structs

struct Rename {
    1: required string Default (go.tag = 'json:"default"')
    2: required string camelCase (go.tag = 'json:"snake_case"')
}

struct Omit {
    1: required string serialized
    2: required string hidden (go.tag = 'json:"-"')
}

struct GoTags {
        1: required string Foo (go.tag = 'json:"-" foo:"bar"')
        2: optional string Bar (go.tag = 'bar:"foo"')
        3: required string FooBar (go.tag = 'json:"foobar,option1,option2" bar:"foo,option1" foo:"foobar"')
        4: required string FooBarWithSpace (go.tag = 'json:"foobarWithSpace" foo:"foo bar foobar barfoo"')
        5: optional string FooBarWithOmitEmpty (go.tag = 'json:"foobarWithOmitEmpty,omitempty"')
        6: required string FooBarWithRequired (go.tag = 'json:"foobarWithRequired,required"')
}

struct NotOmitEmpty {
    1: optional string NotOmitEmptyString (go.tag = 'json:"notOmitEmptyString,!omitempty"')
    2: optional string NotOmitEmptyInt (go.tag = 'json:"notOmitEmptyInt,!omitempty"')
    3: optional string NotOmitEmptyBool (go.tag = 'json:"notOmitEmptyBool,!omitempty"')
    4: optional list<string> NotOmitEmptyList (go.tag = 'json:"notOmitEmptyList,!omitempty"')
    5: optional map<string, string> NotOmitEmptyMap (go.tag = 'json:"notOmitEmptyMap,!omitempty"')
    6: optional list<string> NotOmitEmptyListMixedWithOmitEmpty (go.tag = 'json:"notOmitEmptyListMixedWithOmitEmpty,!omitempty,omitempty"')
    7: optional list<string> NotOmitEmptyListMixedWithOmitEmptyV2 (go.tag = 'json:"notOmitEmptyListMixedWithOmitEmptyV2,omitempty,!omitempty"')
    8: optional string OmitEmptyString (go.tag = 'json:"omitEmptyString,omitempty"') // to test that there can be a mix of fields that do and don't have !omitempty
}

//////////////////////////////////////////////////////////////////////////////
// Default values

struct DefaultsStruct {
    1: required i32 requiredPrimitive = 100
    2: optional i32 optionalPrimitive = 200

    3: required enums.EnumDefault requiredEnum = enums.EnumDefault.Bar
    4: optional enums.EnumDefault optionalEnum = 2

    5: required list<string> requiredList = ["hello", "world"]
    6: optional list<double> optionalList = [1, 2.0, 3]

    7: required Frame requiredStruct = {
        "topLeft": {"x": 1, "y": 2},
        "size": {"width": 100, "height": 200},
    }
    8: optional Edge optionalStruct = {
        "startPoint": {"x": 1, "y": 2},
        "endPoint":   {"x": 3, "y": 4},
    }

    9:  required bool requiredBoolDefaultTrue = true
    10: optional bool optionalBoolDefaultTrue = true

    11: required bool requiredBoolDefaultFalse = false
    12: optional bool optionalBoolDefaultFalse = false
}

//////////////////////////////////////////////////////////////////////////////
// Opt-out of Zap

struct ZapOptOutStruct {
    1: required string name
    2: required string optout (go.nolog)
}

//////////////////////////////////////////////////////////////////////////////
// Field jabels

struct StructLabels {
    // reserved keyword as label
    1: optional bool isRequired (go.label = "required")

    // go.tag's JSON tag takes precedence over go.label
    2: optional string foo (go.label = "bar", go.tag = 'json:"not_bar"')

    // Empty label
    3: optional string qux (go.label = "")

    // All-caps label
    4: optional string quux (go.label = "QUUX")
}

/**
 * Number of seconds since epoch.
 *
 * Deprecated: Use ISOTime instead.
 */
typedef i64 Timestamp  // alias of primitive
typedef string State

typedef i128 UUID  // alias of struct

typedef UUID MyUUID // alias of alias

typedef list<Event> EventGroup  // alias fo collection

struct i128 {
    1: required i64 high
    2: required i64 low
}

struct Event {
    1: required UUID uuid  // required typedef
    2: optional Timestamp time  // optional typedef
}

struct TransitiveTypedefField {
    1: required MyUUID defUUID  // required typedef of alias
}

struct DefaultPrimitiveTypedef {
    1: optional State state = "hello"
}

struct Transition {
    1: required State fromState
    2: required State toState
    3: optional EventGroup events
}

typedef binary PDF  // alias of []byte

typedef set<structs.Frame> FrameGroup

typedef map<structs.Point, structs.Point> PointMap

typedef set<binary> BinarySet

typedef map<structs.Edge, structs.Edge> EdgeMap

typedef map<State, i64> StateMap

typedef enums.EnumWithValues MyEnum

union EmptyUnion {}

union Document {
    1: typedefs.PDF pdf
    2: string plainText
}

/**
 * ArbitraryValue allows constructing complex values without a schema.
 *
 * A value is one of,
 *
 * * Boolean
 * * Integer
 * * String
 * * A list of other values
 * * A dictionary of other values
 */
union ArbitraryValue {
    1: bool boolValue
    2: i64 int64Value
    3: string stringValue
    4: list<ArbitraryValue> listValue
    5: map<string, ArbitraryValue> mapValue
}

typedef string UUID

struct UUIDConflict {
    1: required UUID localUUID
    2: required typedefs.UUID importedUUID
}