// This test file is adapted from all test files in https://github.com/thriftrw/thriftrw-go/tree/183888fb47c3c225e86d634fa3701ce8b84c1914/gen/internal/tests/thrift
include "./other_constants.thrift"
include "./containers.thrift"
include "./enums.thrift"
include "./exceptions.thrift"
include "./structs.thrift"
include "./unions.thrift"
include "./typedefs.thrift"

struct StructCollision {
	1: required bool collisionField
	2: required string collision_field (go.name = "CollisionField2")
}

struct struct_collision {
	1: required bool collisionField
	2: required string collision_field (go.name = "CollisionField2")
} (go.name="StructCollision2")

struct PrimitiveContainers {
    1: optional list<string> ListOrSetOrMap (go.name = "A")
    3: optional set<string>  List_Or_SetOrMap (go.name = "B")
    5: optional map<string, string> ListOrSet_Or_Map (go.name = "C")
}

enum MyEnum {
    X = 123,
    Y = 456,
    Z = 789,
    FooBar,
    foo_bar (go.name="FooBar2"),
}

enum my_enum {
    X = 12,
    Y = 34,
    Z = 56,
} (go.name="MyEnum2")

typedef i64 LittlePotatoe
typedef double little_potatoe (go.name="LittlePotatoe2")

const struct_collision struct_constant = {
	"collisionField": false,
	"collision_field": "false indeed",
}

union UnionCollision {
	1: bool collisionField
	2: string collision_field (go.name = "CollisionField2")
}

union union_collision {
	1: bool collisionField
	2: string collision_field (go.name = "CollisionField2")
} (go.name="UnionCollision2")

struct WithDefault {
	1: required struct_collision pouet = struct_constant
}

struct AccessorNoConflict {
    1: optional string getname
    2: optional string get_name
}

struct AccessorConflict {
    1: optional string name
    2: optional string get_name (go.name = "GetName2")
    3: optional bool is_set_name (go.name = "IsSetName2")
}

const containers.PrimitiveContainers primitiveContainers = {
    "listOfInts": other_constants.listOfInts, // imported constant
    "setOfStrings": ["foo", "bar"],
    "setOfBytes": other_constants.listOfInts, // imported constant with type casting
    "mapOfIntToString": {
        1: "1",
        2: "2",
        3: "3",
    },
    "mapOfStringToBool": {
        "1": 0,
        "2": 1,
        "3": 1,
    }
}

const containers.EnumContainers enumContainers = {
    "listOfEnums": [1, enums.EnumDefault.Foo],
    "setOfEnums": [123, enums.EnumWithValues.Y],
    "mapOfEnums": {
        0: 1,
        enums.EnumWithDuplicateValues.Q: 2,
    },
}

const containers.ContainersOfContainers containersOfContainers = {
    "listOfLists": [[1, 2, 3], [4, 5, 6]],
    "listOfSets": [[1, 2, 3], [4, 5, 6]],
    "listOfMaps": [{1: 2, 3: 4, 5: 6}, {7: 8, 9: 10, 11: 12}],
    "setOfSets": [["1", "2", "3"], ["4", "5", "6"]],
    "setOfLists": [["1", "2", "3"], ["4", "5", "6"]],
    "setOfMaps": [
        {"1": "2", "3": "4", "5": "6"},
        {"7": "8", "9": "10", "11": "12"},
    ],
    "mapOfMapToInt": {
        {"1": 1, "2": 2, "3": 3}: 100,
        {"4": 4, "5": 5, "6": 6}: 200,
    },
    "mapOfListToSet": {
        // more type casting
        other_constants.listOfInts: other_constants.listOfInts,
        [4, 5, 6]: [4, 5, 6],
    },
    "mapOfSetToListOfDouble": {
        [1, 2, 3]: [1.2, 3.4],
        [4, 5, 6]: [5.6, 7.8],
    },
}

const enums.StructWithOptionalEnum structWithOptionalEnum = {
    "e": enums.EnumDefault.Baz
}

const exceptions.EmptyException emptyException = {}

const structs.Graph graph = {
    "edges": [
        {"startPoint": other_constants.some_point, "endPoint": {"x": 3, "y": 4}},
        {"startPoint": {"x": 5, "y": 6}, "endPoint": {"x": 7, "y": 8}},
    ]
}

const structs.Node lastNode = {"value": 3}
const structs.Node node = {
    "value": 1,
    "tail": {"value": 2, "tail": lastNode},
}

const unions.ArbitraryValue arbitraryValue = {
    "listValue": [
        {"boolValue": 1},
        {"int64Value": 2},
        {"stringValue": "hello"},
        {"mapValue": {"foo": {"stringValue": "bar"}}},
    ],
}
// TODO: union validation for constants?

const typedefs.i128 i128 = uuid
const typedefs.UUID uuid = {"high": 1234, "low": 5678}

/** Timestamp at which time began. */
const typedefs.Timestamp beginningOfTime = 0

/**
 * An example frame group.
 *
 * Contains two frames.
 */
const typedefs.FrameGroup frameGroup = [
    {
        "topLeft": {"x": 1, "y": 2},
        "size": {"width": 100, "height": 200},
    }
    {
        "topLeft": {"x": 3, "y": 4},
        "size": {"width": 300, "height": 400},
    },
]

const typedefs.MyEnum myEnum = enums.EnumWithValues.Y

const enums.RecordType NAME = enums.RecordType.NAME
const enums.RecordType HOME = enums.RecordType.HOME_ADDRESS
const enums.RecordType WORK_ADDRESS = enums.RecordType.WORK_ADDRESS

const enums.lowerCaseEnum lower = enums.lowerCaseEnum.items

struct PrimitiveContainers {
    1: optional list<binary> listOfBinary
    2: optional list<i64> listOfInts
    3: optional set<string> setOfStrings
    4: optional set<byte> setOfBytes
    5: optional map<i32, string> mapOfIntToString
    6: optional map<string, bool> mapOfStringToBool
}

struct PrimitiveContainersRequired {
    1: required list<string> listOfStrings
    2: required set<i32> setOfInts
    3: required map<i64, double> mapOfIntsToDoubles
}

struct EnumContainers {
    1: optional list<enums.EnumDefault> listOfEnums
    2: optional set<enums.EnumWithValues> setOfEnums
    3: optional map<enums.EnumWithDuplicateValues, i32> mapOfEnums
}

struct ContainersOfContainers {
    1: optional list<list<i32>> listOfLists;
    2: optional list<set<i32>> listOfSets;
    3: optional list<map<i32, i32>> listOfMaps;

    4: optional set<set<string>> setOfSets;
    5: optional set<list<string>> setOfLists;
    6: optional set<map<string, string>> setOfMaps;

    7: optional map<map<string, i32>, i64> mapOfMapToInt;
    8: optional map<list<i32>, set<i64>> mapOfListToSet;
    9: optional map<set<i32>, list<double>> mapOfSetToListOfDouble;
}

struct MapOfBinaryAndString {
    1: optional map<binary, string> binaryToString;
    2: optional map<string, binary> stringToBinary;
}

struct ListOfRequiredPrimitives {
    1: required list<string> listOfStrings
}

struct ListOfOptionalPrimitives {
    1: optional list<string> listOfStrings
}

struct ListOfConflictingEnums {
    1: required list<enum_conflict.RecordType> records
    2: required list<enums.RecordType> otherRecords
}

struct ListOfConflictingUUIDs {
    1: required list<typedefs.UUID> uuids
    2: required list<uuid_conflict.UUID> otherUUIDs
}

enum EnumMarshalStrict {
    Foo, Bar, Baz, Bat
}

enum RecordType {
    Name, Email
}

const RecordType defaultRecordType = RecordType.Name

const enums.RecordType defaultOtherRecordType = enums.RecordType.NAME

struct Records {
    1: optional RecordType recordType = defaultRecordType
    2: optional enums.RecordType otherRecordType = defaultOtherRecordType
}

enum EmptyEnum {}

enum EnumDefault {
    Foo, Bar, Baz
}

enum EnumWithValues {
    X = 123,
    Y = 456,
    Z = 789,
}

enum EnumWithDuplicateValues {
    P, // 0
    Q = -1,
    R, // 0
}

// enum with item names conflicting with those of another enum
enum EnumWithDuplicateName {
    A, B, C, P, Q, R, X, Y, Z
}

// Enum treated as optional inside a struct
struct StructWithOptionalEnum {
    1: optional EnumDefault e
}

/**
 * Kinds of records stored in the database.
 */
enum RecordType {
  /** Name of the user. */
  NAME,

  /**
   * Home address of the user.
   *
   * This record is always present.
   */
  HOME_ADDRESS,

  /**
   * Home address of the user.
   *
   * This record may not be present.
   */
  WORK_ADDRESS
}

enum lowerCaseEnum {
    containing, lower_case, items
}

// EnumWithLabel use label name in serialization/deserialization
enum EnumWithLabel {
    USERNAME (go.label = "surname"),
    PASSWORD (go.label = "hashed_password"),
    SALT (go.label = ""),
    SUGAR (go.label),
    relay (go.label = "RELAY")
    NAIVE4_N1 (go.label = "function")

}

// collision with RecordType_Values() function.
enum RecordType_Values { FOO, BAR }

exception EmptyException {}

/**
 * Raised when something doesn't exist.
 */
exception DoesNotExistException {
    /** Key that was missing. */
    1: required string key
    2: optional string Error (go.name="Error2")
}

exception Does_Not_Exist_Exception_Collision {
 /** Key that was missing. */
    1: required string key
    2: optional string Error (go.name="Error2")
} (go.name="DoesNotExistException2")

struct DocumentStruct {
 1: required non_hyphenated.Second second
}

struct DocumentStructure {
 1: required non_hyphenated.Second r2
}