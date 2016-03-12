//line sql.y:6
package sqlparser

import __yyfmt__ "fmt"

//line sql.y:6
import "bytes"

func SetParseTree(yylex interface{}, stmt Statement) {
	yylex.(*Tokenizer).ParseTree = stmt
}

func SetAllowComments(yylex interface{}, allow bool) {
	yylex.(*Tokenizer).AllowComments = allow
}

func ForceEOF(yylex interface{}) {
	yylex.(*Tokenizer).ForceEOF = true
}

var (
	SHARE        = []byte("share")
	MODE         = []byte("mode")
	IF_BYTES     = []byte("if")
	VALUES_BYTES = []byte("values")
)

//line sql.y:31
type yySymType struct {
	yys         int
	empty       struct{}
	statement   Statement
	statements  Statements
	selStmt     SelectStatement
	byt         byte
	bytes       []byte
	bytes2      [][]byte
	str         string
	selectExprs SelectExprs
	selectExpr  SelectExpr
	columns     Columns
	colName     *ColName
	tableExprs  TableExprs
	tableExpr   TableExpr
	smTableExpr SimpleTableExpr
	tableName   *TableName
	indexHints  *IndexHints
	expr        Expr
	boolExpr    BoolExpr
	valExpr     ValExpr
	colTuple    ColTuple
	valExprs    ValExprs
	values      Values
	rowTuple    RowTuple
	subquery    *Subquery
	caseExpr    *CaseExpr
	whens       []*When
	when        *When
	orderBy     OrderBy
	order       *Order
	limit       *Limit
	insRows     InsertRows
	updateExprs UpdateExprs
	updateExpr  *UpdateExpr
}

const LEX_ERROR = 57346
const SELECT = 57347
const INSERT = 57348
const UPDATE = 57349
const DELETE = 57350
const FROM = 57351
const WHERE = 57352
const GROUP = 57353
const HAVING = 57354
const ORDER = 57355
const BY = 57356
const LIMIT = 57357
const FOR = 57358
const ALL = 57359
const DISTINCT = 57360
const AS = 57361
const EXISTS = 57362
const IN = 57363
const IS = 57364
const LIKE = 57365
const BETWEEN = 57366
const NULL = 57367
const ASC = 57368
const DESC = 57369
const VALUES = 57370
const INTO = 57371
const DUPLICATE = 57372
const KEY = 57373
const DEFAULT = 57374
const SET = 57375
const LOCK = 57376
const KEYRANGE = 57377
const ID = 57378
const STRING = 57379
const NUMBER = 57380
const VALUE_ARG = 57381
const LIST_ARG = 57382
const COMMENT = 57383
const LE = 57384
const GE = 57385
const NE = 57386
const NULL_SAFE_EQUAL = 57387
const UNION = 57388
const MINUS = 57389
const EXCEPT = 57390
const INTERSECT = 57391
const JOIN = 57392
const STRAIGHT_JOIN = 57393
const LEFT = 57394
const RIGHT = 57395
const INNER = 57396
const OUTER = 57397
const CROSS = 57398
const NATURAL = 57399
const USE = 57400
const FORCE = 57401
const ON = 57402
const OR = 57403
const AND = 57404
const NOT = 57405
const UNARY = 57406
const CASE = 57407
const WHEN = 57408
const THEN = 57409
const ELSE = 57410
const END = 57411
const CREATE = 57412
const ALTER = 57413
const DROP = 57414
const RENAME = 57415
const ANALYZE = 57416
const TABLE = 57417
const INDEX = 57418
const VIEW = 57419
const TO = 57420
const IGNORE = 57421
const IF = 57422
const UNIQUE = 57423
const USING = 57424
const SHOW = 57425
const DESCRIBE = 57426
const EXPLAIN = 57427

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"LEX_ERROR",
	"SELECT",
	"INSERT",
	"UPDATE",
	"DELETE",
	"FROM",
	"WHERE",
	"GROUP",
	"HAVING",
	"ORDER",
	"BY",
	"LIMIT",
	"FOR",
	"ALL",
	"DISTINCT",
	"AS",
	"EXISTS",
	"IN",
	"IS",
	"LIKE",
	"BETWEEN",
	"NULL",
	"ASC",
	"DESC",
	"VALUES",
	"INTO",
	"DUPLICATE",
	"KEY",
	"DEFAULT",
	"SET",
	"LOCK",
	"KEYRANGE",
	"ID",
	"STRING",
	"NUMBER",
	"VALUE_ARG",
	"LIST_ARG",
	"COMMENT",
	"LE",
	"GE",
	"NE",
	"NULL_SAFE_EQUAL",
	"'('",
	"'='",
	"'<'",
	"'>'",
	"'~'",
	"UNION",
	"MINUS",
	"EXCEPT",
	"INTERSECT",
	"','",
	"JOIN",
	"STRAIGHT_JOIN",
	"LEFT",
	"RIGHT",
	"INNER",
	"OUTER",
	"CROSS",
	"NATURAL",
	"USE",
	"FORCE",
	"ON",
	"OR",
	"AND",
	"NOT",
	"'&'",
	"'|'",
	"'^'",
	"'+'",
	"'-'",
	"'*'",
	"'/'",
	"'%'",
	"'.'",
	"UNARY",
	"CASE",
	"WHEN",
	"THEN",
	"ELSE",
	"END",
	"CREATE",
	"ALTER",
	"DROP",
	"RENAME",
	"ANALYZE",
	"TABLE",
	"INDEX",
	"VIEW",
	"TO",
	"IGNORE",
	"IF",
	"UNIQUE",
	"USING",
	"SHOW",
	"DESCRIBE",
	"EXPLAIN",
	"';'",
	"')'",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyNprod = 207
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 587

var yyAct = [...]int{

	100, 336, 68, 299, 167, 369, 164, 253, 97, 91,
	98, 55, 203, 291, 96, 244, 183, 214, 87, 349,
	69, 166, 5, 86, 378, 7, 378, 191, 378, 109,
	82, 264, 265, 266, 267, 268, 135, 269, 270, 141,
	140, 56, 57, 71, 74, 108, 76, 260, 114, 79,
	235, 70, 297, 83, 7, 60, 72, 105, 106, 107,
	34, 35, 36, 37, 78, 135, 170, 135, 347, 92,
	112, 380, 235, 379, 233, 377, 125, 320, 43, 48,
	45, 49, 129, 325, 46, 133, 316, 318, 126, 346,
	137, 128, 345, 110, 111, 75, 54, 322, 50, 296,
	115, 275, 168, 327, 163, 165, 169, 245, 51, 52,
	53, 234, 285, 139, 283, 113, 317, 122, 71, 236,
	181, 71, 177, 187, 118, 245, 70, 289, 186, 70,
	140, 141, 140, 173, 124, 188, 153, 154, 155, 209,
	187, 185, 92, 141, 140, 201, 329, 223, 342, 213,
	211, 212, 221, 222, 208, 225, 226, 227, 228, 229,
	230, 231, 232, 323, 207, 148, 149, 150, 151, 152,
	153, 154, 155, 216, 77, 292, 344, 237, 92, 92,
	71, 71, 249, 151, 152, 153, 154, 155, 70, 251,
	242, 224, 257, 197, 255, 239, 241, 65, 252, 120,
	248, 256, 132, 310, 258, 343, 308, 314, 311, 313,
	292, 309, 195, 210, 312, 134, 198, 120, 274, 261,
	276, 235, 237, 354, 184, 331, 278, 279, 286, 18,
	19, 20, 21, 184, 81, 207, 277, 264, 265, 266,
	267, 268, 282, 269, 270, 121, 206, 92, 216, 18,
	18, 34, 35, 36, 37, 290, 205, 22, 294, 288,
	298, 135, 116, 295, 284, 119, 194, 196, 193, 262,
	108, 217, 179, 114, 364, 306, 307, 215, 120, 363,
	206, 72, 105, 106, 107, 180, 362, 324, 84, 170,
	205, 170, 207, 207, 175, 112, 328, 174, 71, 172,
	171, 326, 333, 108, 58, 77, 332, 334, 337, 23,
	24, 26, 25, 27, 338, 105, 106, 107, 110, 111,
	273, 138, 28, 29, 30, 115, 72, 321, 319, 303,
	348, 302, 200, 199, 182, 130, 350, 272, 77, 127,
	113, 123, 352, 66, 85, 80, 360, 358, 237, 375,
	359, 117, 361, 6, 351, 330, 367, 18, 64, 281,
	366, 337, 368, 370, 370, 370, 71, 376, 373, 371,
	372, 240, 382, 103, 70, 189, 131, 63, 108, 383,
	247, 114, 61, 384, 300, 385, 59, 254, 104, 90,
	105, 106, 107, 353, 218, 341, 219, 220, 301, 95,
	340, 305, 184, 112, 67, 381, 365, 18, 148, 149,
	150, 151, 152, 153, 154, 155, 4, 356, 357, 190,
	44, 259, 94, 192, 103, 47, 110, 111, 88, 108,
	73, 250, 114, 115, 178, 374, 355, 335, 339, 104,
	90, 105, 106, 107, 304, 287, 176, 243, 113, 102,
	95, 99, 101, 293, 112, 238, 246, 142, 93, 315,
	18, 148, 149, 150, 151, 152, 153, 154, 155, 204,
	263, 202, 89, 94, 271, 103, 136, 110, 111, 88,
	108, 62, 33, 114, 115, 31, 17, 16, 15, 14,
	104, 72, 105, 106, 107, 103, 3, 13, 12, 113,
	108, 95, 32, 114, 11, 112, 10, 9, 8, 2,
	104, 72, 105, 106, 107, 38, 39, 40, 41, 42,
	1, 95, 0, 0, 94, 112, 0, 0, 110, 111,
	143, 147, 145, 146, 280, 115, 148, 149, 150, 151,
	152, 153, 154, 155, 94, 0, 0, 0, 110, 111,
	113, 159, 160, 161, 162, 115, 156, 157, 158, 148,
	149, 150, 151, 152, 153, 154, 155, 0, 0, 0,
	113, 0, 0, 0, 0, 0, 0, 0, 144, 148,
	149, 150, 151, 152, 153, 154, 155,
}
var yyPact = [...]int{

	-1000, -1000, -79, 224, -1000, -1000, -1000, 200, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -12, -13, 8, 18, 6, -1000, -1000,
	-1000, 263, 224, 402, 365, -1000, -1000, -1000, 359, 329,
	307, 395, 290, -51, 4, 269, -1000, -26, 269, -1000,
	309, -65, 269, -65, 308, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, 404, -1000, 307, 318, 46, 307, 162, -1000,
	198, -1000, 39, 305, 65, 269, -1000, -1000, 303, -1000,
	-11, 299, 356, 136, 269, -1000, 206, -1000, -1000, 302,
	35, 76, 509, -1000, 475, 455, -1000, -1000, -1000, 20,
	254, 253, -1000, 251, 248, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, 20, 239, 290, 298, 392,
	290, 20, 269, -1000, 355, -70, -1000, 180, -1000, 297,
	-1000, -1000, 296, -1000, 210, 404, -1000, -1000, 269, 138,
	475, 475, 20, 231, 373, 20, 20, 122, 20, 20,
	20, 20, 20, 20, 20, 20, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, 509, -28, 9, 17, 509, -1000,
	245, 353, 404, -1000, 402, 278, 26, 489, 352, 290,
	290, 223, -1000, 374, 475, -1000, 489, -1000, -1000, -1000,
	135, 269, -1000, -46, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, 214, 181, 301, 244, 23, -1000, -1000, -1000,
	-1000, -1000, 62, 489, -1000, 245, -1000, -1000, 231, 20,
	20, 489, 466, -1000, 334, 110, 110, 110, 61, 61,
	-1000, -1000, -1000, -1000, -1000, 20, -1000, 489, -1000, 12,
	404, 10, 173, 44, -1000, 475, 109, 243, 200, 144,
	-3, -1000, 374, 369, 384, 76, 295, -1000, -1000, 293,
	-1000, 390, 210, 210, -1000, -1000, 150, 147, 158, 153,
	151, 22, -1000, 292, -25, 291, -5, -1000, 489, 95,
	20, -1000, 489, -1000, -19, -1000, 278, 19, -1000, 20,
	64, -1000, 325, 170, -1000, -1000, -1000, 290, 369, -1000,
	20, 20, -1000, -1000, 388, 381, 181, 82, -1000, 149,
	-1000, 120, -1000, -1000, -1000, -1000, 1, -2, -23, -1000,
	-1000, -1000, -1000, 20, 489, -1000, -83, -1000, 489, 20,
	323, 243, -1000, -1000, 338, 168, -1000, 391, -1000, 374,
	475, 20, 475, -1000, -1000, 240, 233, 228, 489, -1000,
	489, 399, -1000, 20, 20, -1000, -1000, -1000, 369, 76,
	166, 76, 269, 269, 269, 290, 489, -1000, 333, -27,
	-1000, -29, -31, 162, -1000, 398, 351, -1000, 269, -1000,
	-1000, -1000, 269, -1000, 269, -1000,
}
var yyPgo = [...]int{

	0, 520, 353, 509, 21, 508, 507, 506, 504, 498,
	497, 489, 488, 487, 486, 496, 485, 482, 481, 23,
	18, 476, 474, 472, 471, 12, 470, 469, 197, 459,
	5, 16, 9, 458, 457, 456, 14, 6, 17, 4,
	453, 10, 452, 29, 451, 8, 449, 447, 15, 446,
	445, 444, 438, 7, 437, 1, 436, 3, 435, 434,
	431, 13, 2, 20, 234, 430, 425, 423, 421, 420,
	419, 0, 11, 416,
}
var yyR1 = [...]int{

	0, 1, 3, 3, 3, 2, 2, 2, 2, 2,
	2, 2, 2, 2, 2, 2, 4, 4, 4, 5,
	5, 6, 7, 8, 9, 9, 9, 10, 10, 10,
	11, 12, 12, 12, 13, 14, 14, 14, 73, 15,
	16, 16, 17, 17, 17, 17, 17, 18, 18, 19,
	19, 20, 20, 20, 23, 23, 21, 21, 21, 24,
	24, 25, 25, 25, 25, 22, 22, 22, 26, 26,
	26, 26, 26, 26, 26, 26, 26, 27, 27, 27,
	28, 28, 29, 29, 29, 29, 30, 30, 31, 31,
	32, 32, 32, 32, 32, 33, 33, 33, 33, 33,
	33, 33, 33, 33, 33, 33, 34, 34, 34, 34,
	34, 34, 34, 38, 38, 38, 43, 39, 39, 37,
	37, 37, 37, 37, 37, 37, 37, 37, 37, 37,
	37, 37, 37, 37, 37, 37, 42, 42, 44, 44,
	44, 46, 49, 49, 47, 47, 48, 50, 50, 45,
	45, 36, 36, 36, 36, 51, 51, 52, 52, 53,
	53, 54, 54, 55, 56, 56, 56, 57, 57, 57,
	58, 58, 58, 59, 59, 60, 60, 61, 61, 35,
	35, 40, 40, 41, 41, 62, 62, 63, 64, 64,
	65, 65, 66, 66, 67, 67, 67, 67, 67, 68,
	68, 69, 69, 70, 70, 71, 72,
}
var yyR2 = [...]int{

	0, 1, 2, 4, 3, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 12, 4, 3, 7,
	7, 8, 7, 3, 5, 8, 4, 6, 7, 4,
	5, 4, 5, 5, 3, 2, 2, 2, 0, 2,
	0, 2, 1, 2, 1, 1, 1, 0, 1, 1,
	3, 1, 2, 3, 1, 1, 0, 1, 2, 1,
	3, 3, 3, 3, 5, 0, 1, 2, 1, 1,
	2, 3, 2, 3, 2, 2, 2, 1, 3, 1,
	1, 3, 0, 5, 5, 5, 1, 3, 0, 2,
	1, 3, 3, 2, 3, 3, 3, 4, 3, 4,
	5, 6, 3, 4, 2, 6, 1, 1, 1, 1,
	1, 1, 1, 3, 1, 1, 3, 1, 3, 1,
	1, 1, 3, 3, 3, 3, 3, 3, 3, 3,
	2, 3, 4, 5, 4, 1, 1, 1, 1, 1,
	1, 5, 0, 1, 1, 2, 4, 0, 2, 1,
	3, 1, 1, 1, 1, 0, 3, 0, 2, 0,
	3, 1, 3, 2, 0, 1, 1, 0, 2, 4,
	0, 2, 4, 0, 3, 1, 3, 0, 5, 2,
	1, 1, 3, 3, 1, 1, 3, 3, 0, 2,
	0, 3, 0, 1, 1, 1, 1, 1, 1, 0,
	1, 0, 1, 0, 2, 1, 0,
}
var yyChk = [...]int{

	-1000, -1, -3, -15, -73, 101, -2, -4, -5, -6,
	-7, -8, -9, -10, -11, -12, -13, -14, 5, 6,
	7, 8, 33, 85, 86, 88, 87, 89, 98, 99,
	100, -16, -15, -17, 51, 52, 53, 54, -15, -15,
	-15, -15, -15, 90, -69, 92, 96, -66, 92, 94,
	90, 90, 91, 92, 90, -72, -72, -72, 41, -2,
	-4, 17, -18, 18, 29, -28, 36, 9, -62, -63,
	-45, -71, 36, -65, 95, 91, -71, 36, 90, -71,
	36, -64, 95, -71, -64, 36, -19, -20, 75, -23,
	36, -32, -37, -33, 69, 46, -36, -45, -41, -44,
	-71, -42, -46, 20, 35, 37, 38, 39, 25, -43,
	73, 74, 50, 95, 28, 80, -28, 33, 78, -28,
	55, 47, 78, 36, 69, -71, -72, 36, -72, 93,
	36, 20, 66, -71, 9, 55, -21, -71, 19, 78,
	68, 67, -34, 21, 69, 23, 24, 22, 70, 71,
	72, 73, 74, 75, 76, 77, 47, 48, 49, 42,
	43, 44, 45, -32, -37, -32, -4, -39, -37, -37,
	46, 46, 46, -43, 46, 46, -49, -37, -59, 33,
	46, -62, 36, -31, 10, -63, -37, -71, -72, 20,
	-70, 97, -67, 88, 86, 32, 87, 13, 36, 36,
	36, -72, -24, -25, -27, 46, 36, -43, -20, -71,
	75, -32, -32, -37, -38, 46, -43, 40, 21, 23,
	24, -37, -37, 25, 69, -37, -37, -37, -37, -37,
	-37, -37, -37, 102, 102, 55, 102, -37, 102, -19,
	18, -19, -36, -47, -48, 81, -35, 28, -4, -62,
	-60, -45, -31, -53, 13, -32, 66, -71, -72, -68,
	93, -31, 55, -26, 56, 57, 58, 59, 60, 62,
	63, -22, 36, 19, -25, 78, -39, -38, -37, -37,
	68, 25, -37, 102, -19, 102, 55, -50, -48, 83,
	-32, -61, 66, -40, -41, -61, 102, 55, -53, -57,
	15, 14, 36, 36, -51, 11, -25, -25, 56, 61,
	56, 61, 56, 56, 56, -29, 64, 94, 65, 36,
	102, 36, 102, 68, -37, 102, -36, 84, -37, 82,
	30, 55, -45, -57, -37, -54, -55, -37, -72, -52,
	12, 14, 66, 56, 56, 91, 91, 91, -37, 102,
	-37, 31, -41, 55, 55, -56, 26, 27, -53, -32,
	-39, -32, 46, 46, 46, 7, -37, -55, -57, -30,
	-71, -30, -30, -62, -58, 16, 34, 102, 55, 102,
	102, 7, 21, -71, -71, -71,
}
var yyDef = [...]int{

	38, -2, 1, 0, 40, 38, 2, 5, 6, 7,
	8, 9, 10, 11, 12, 13, 14, 15, 38, 38,
	38, 38, 38, 201, 192, 0, 0, 0, 206, 206,
	206, 39, 4, 0, 42, 44, 45, 46, 47, 0,
	0, 0, 0, 190, 0, 0, 202, 0, 0, 193,
	0, 188, 0, 188, 0, 35, 36, 37, 41, 3,
	18, 43, 0, 48, 0, 0, 80, 0, 23, 185,
	0, 149, 205, 0, 0, 0, 206, 205, 0, 206,
	0, 0, 0, 0, 0, 34, 17, 49, 51, 56,
	205, 54, 55, 90, 0, 0, 119, 120, 121, 0,
	149, 0, 135, 0, 0, 151, 152, 153, 154, 184,
	138, 139, 140, 136, 137, 142, 173, 0, 0, 88,
	0, 0, 0, 206, 0, 203, 26, 0, 29, 0,
	31, 189, 0, 206, 0, 0, 52, 57, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 106, 107, 108, 109,
	110, 111, 112, 93, 0, 0, 0, 0, 117, 130,
	0, 0, 0, 104, 0, 0, 0, 143, 0, 0,
	0, 88, 81, 159, 0, 186, 187, 150, 24, 191,
	0, 0, 206, 199, 194, 195, 196, 197, 198, 30,
	32, 33, 88, 59, 65, 0, 77, 79, 50, 58,
	53, 91, 92, 95, 96, 0, 114, 115, 0, 0,
	0, 98, 0, 102, 0, 122, 123, 124, 125, 126,
	127, 128, 129, 94, 116, 0, 183, 117, 131, 0,
	0, 0, 0, 147, 144, 0, 177, 0, 180, 177,
	0, 175, 159, 167, 0, 89, 0, 204, 27, 0,
	200, 155, 0, 0, 68, 69, 0, 0, 0, 0,
	0, 82, 66, 0, 0, 0, 0, 97, 99, 0,
	0, 103, 118, 132, 0, 134, 0, 0, 145, 0,
	0, 19, 0, 179, 181, 20, 174, 0, 167, 22,
	0, 0, 206, 28, 157, 0, 60, 63, 70, 0,
	72, 0, 74, 75, 76, 61, 0, 0, 0, 67,
	62, 78, 113, 0, 100, 133, 0, 141, 148, 0,
	0, 0, 176, 21, 168, 160, 161, 164, 25, 159,
	0, 0, 0, 71, 73, 0, 0, 0, 101, 105,
	146, 0, 182, 0, 0, 163, 165, 166, 167, 158,
	156, 64, 0, 0, 0, 0, 169, 162, 170, 0,
	86, 0, 0, 178, 16, 0, 0, 83, 0, 84,
	85, 171, 0, 87, 0, 172,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 77, 70, 3,
	46, 102, 75, 73, 55, 74, 78, 76, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 101,
	48, 47, 49, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 72, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 71, 3, 50,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39, 40, 41,
	42, 43, 44, 45, 51, 52, 53, 54, 56, 57,
	58, 59, 60, 61, 62, 63, 64, 65, 66, 67,
	68, 69, 79, 80, 81, 82, 83, 84, 85, 86,
	87, 88, 89, 90, 91, 92, 93, 94, 95, 96,
	97, 98, 99, 100,
}
var yyTok3 = [...]int{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lval  yySymType
	stack [yyInitialStackSize]yySymType
	char  int
}

func (p *yyParserImpl) Lookahead() int {
	return p.char
}

func yyNewParser() yyParser {
	return &yyParserImpl{}
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := yyPact[state]
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && yyChk[yyAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || yyExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := yyExca[i]
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		token = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = yyTok3[i+0]
		if token == char {
			token = yyTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := yyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yyrcvr.char = -1
	yytoken := -1 // yyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yyrcvr.char = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yyrcvr.char < 0 {
		yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yytoken { /* valid shift */
		yyrcvr.char = -1
		yytoken = -1
		yyVAL = yyrcvr.lval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yyrcvr.char < 0 {
			yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yyrcvr.char = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:153
		{
			SetParseTree(yylex, yyDollar[1].statements)
		}
	case 2:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:159
		{
			comment := Comments(yyDollar[1].bytes2)
			if comment.IsEmpty() {
				yyVAL.statements = Statements([]Statement{yyDollar[2].statement})
			} else {
				yyVAL.statements = Statements([]Statement{comment, yyDollar[2].statement})
			}
		}
	case 3:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line sql.y:168
		{
			comment := Comments(yyDollar[3].bytes2)
			if comment.IsEmpty() {
				yyVAL.statements = Statements(append(yyDollar[1].statements, yyDollar[4].statement))
			} else {
				yyVAL.statements = Statements(append(yyDollar[1].statements, comment, yyDollar[4].statement))
			}
		}
	case 4:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:177
		{
			comment := Comments(yyDollar[3].bytes2)
			if comment.IsEmpty() {
				yyVAL.statements = yyDollar[1].statements
			} else {
				yyVAL.statements = Statements(append(yyDollar[1].statements, comment))
			}
		}
	case 5:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:188
		{
			yyVAL.statement = yyDollar[1].selStmt
		}
	case 16:
		yyDollar = yyS[yypt-12 : yypt+1]
		//line sql.y:204
		{
			yyVAL.selStmt = &Select{Comments: Comments(yyDollar[2].bytes2), Distinct: yyDollar[3].str, SelectExprs: yyDollar[4].selectExprs, From: yyDollar[6].tableExprs, Where: NewWhere(AST_WHERE, yyDollar[7].boolExpr), GroupBy: GroupBy(yyDollar[8].valExprs), Having: NewWhere(AST_HAVING, yyDollar[9].boolExpr), OrderBy: yyDollar[10].orderBy, Limit: yyDollar[11].limit, Lock: yyDollar[12].str}
		}
	case 17:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line sql.y:208
		{
			yyVAL.selStmt = &Select{Comments: Comments(yyDollar[2].bytes2), Distinct: yyDollar[3].str, SelectExprs: yyDollar[4].selectExprs}
		}
	case 18:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:212
		{
			yyVAL.selStmt = &Union{Type: yyDollar[2].str, Left: yyDollar[1].selStmt, Right: yyDollar[3].selStmt}
		}
	case 19:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line sql.y:218
		{
			yyVAL.statement = &Insert{Comments: Comments(yyDollar[2].bytes2), Table: yyDollar[4].tableName, Columns: yyDollar[5].columns, Rows: yyDollar[6].insRows, OnDup: OnDup(yyDollar[7].updateExprs)}
		}
	case 20:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line sql.y:222
		{
			cols := make(Columns, 0, len(yyDollar[6].updateExprs))
			vals := make(ValTuple, 0, len(yyDollar[6].updateExprs))
			for _, col := range yyDollar[6].updateExprs {
				cols = append(cols, &NonStarExpr{Expr: col.Name})
				vals = append(vals, col.Expr)
			}
			yyVAL.statement = &Insert{Comments: Comments(yyDollar[2].bytes2), Table: yyDollar[4].tableName, Columns: cols, Rows: Values{vals}, OnDup: OnDup(yyDollar[7].updateExprs)}
		}
	case 21:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line sql.y:234
		{
			yyVAL.statement = &Update{Comments: Comments(yyDollar[2].bytes2), Table: yyDollar[3].tableName, Exprs: yyDollar[5].updateExprs, Where: NewWhere(AST_WHERE, yyDollar[6].boolExpr), OrderBy: yyDollar[7].orderBy, Limit: yyDollar[8].limit}
		}
	case 22:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line sql.y:240
		{
			yyVAL.statement = &Delete{Comments: Comments(yyDollar[2].bytes2), Table: yyDollar[4].tableName, Where: NewWhere(AST_WHERE, yyDollar[5].boolExpr), OrderBy: yyDollar[6].orderBy, Limit: yyDollar[7].limit}
		}
	case 23:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:246
		{
			yyVAL.statement = &Set{Comments: Comments(yyDollar[2].bytes2), Exprs: yyDollar[3].updateExprs}
		}
	case 24:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line sql.y:252
		{
			yyVAL.statement = &DDL{Action: AST_CREATE, NewName: yyDollar[4].bytes}
		}
	case 25:
		yyDollar = yyS[yypt-8 : yypt+1]
		//line sql.y:256
		{
			// Change this to an alter statement
			yyVAL.statement = &DDL{Action: AST_ALTER, Table: yyDollar[7].bytes, NewName: yyDollar[7].bytes}
		}
	case 26:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line sql.y:261
		{
			yyVAL.statement = &DDL{Action: AST_CREATE, NewName: yyDollar[3].bytes}
		}
	case 27:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line sql.y:267
		{
			yyVAL.statement = &DDL{Action: AST_ALTER, Table: yyDollar[4].bytes, NewName: yyDollar[4].bytes}
		}
	case 28:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line sql.y:271
		{
			// Change this to a rename statement
			yyVAL.statement = &DDL{Action: AST_RENAME, Table: yyDollar[4].bytes, NewName: yyDollar[7].bytes}
		}
	case 29:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line sql.y:276
		{
			yyVAL.statement = &DDL{Action: AST_ALTER, Table: yyDollar[3].bytes, NewName: yyDollar[3].bytes}
		}
	case 30:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line sql.y:282
		{
			yyVAL.statement = &DDL{Action: AST_RENAME, Table: yyDollar[3].bytes, NewName: yyDollar[5].bytes}
		}
	case 31:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line sql.y:288
		{
			yyVAL.statement = &DDL{Action: AST_DROP, Table: yyDollar[4].bytes}
		}
	case 32:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line sql.y:292
		{
			// Change this to an alter statement
			yyVAL.statement = &DDL{Action: AST_ALTER, Table: yyDollar[5].bytes, NewName: yyDollar[5].bytes}
		}
	case 33:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line sql.y:297
		{
			yyVAL.statement = &DDL{Action: AST_DROP, Table: yyDollar[4].bytes}
		}
	case 34:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:303
		{
			yyVAL.statement = &DDL{Action: AST_ALTER, Table: yyDollar[3].bytes, NewName: yyDollar[3].bytes}
		}
	case 35:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:309
		{
			yyVAL.statement = &Other{}
		}
	case 36:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:313
		{
			yyVAL.statement = &Other{}
		}
	case 37:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:317
		{
			yyVAL.statement = &Other{}
		}
	case 38:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:322
		{
			SetAllowComments(yylex, true)
		}
	case 39:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:326
		{
			yyVAL.bytes2 = yyDollar[2].bytes2
			SetAllowComments(yylex, false)
		}
	case 40:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:332
		{
			yyVAL.bytes2 = nil
		}
	case 41:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:336
		{
			yyVAL.bytes2 = append(yyDollar[1].bytes2, yyDollar[2].bytes)
		}
	case 42:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:342
		{
			yyVAL.str = AST_UNION
		}
	case 43:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:346
		{
			yyVAL.str = AST_UNION_ALL
		}
	case 44:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:350
		{
			yyVAL.str = AST_SET_MINUS
		}
	case 45:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:354
		{
			yyVAL.str = AST_EXCEPT
		}
	case 46:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:358
		{
			yyVAL.str = AST_INTERSECT
		}
	case 47:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:363
		{
			yyVAL.str = ""
		}
	case 48:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:367
		{
			yyVAL.str = AST_DISTINCT
		}
	case 49:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:373
		{
			yyVAL.selectExprs = SelectExprs{yyDollar[1].selectExpr}
		}
	case 50:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:377
		{
			yyVAL.selectExprs = append(yyVAL.selectExprs, yyDollar[3].selectExpr)
		}
	case 51:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:383
		{
			yyVAL.selectExpr = &StarExpr{}
		}
	case 52:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:387
		{
			yyVAL.selectExpr = &NonStarExpr{Expr: yyDollar[1].expr, As: yyDollar[2].bytes}
		}
	case 53:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:391
		{
			yyVAL.selectExpr = &StarExpr{TableName: yyDollar[1].bytes}
		}
	case 54:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:397
		{
			yyVAL.expr = yyDollar[1].boolExpr
		}
	case 55:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:401
		{
			yyVAL.expr = yyDollar[1].valExpr
		}
	case 56:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:406
		{
			yyVAL.bytes = nil
		}
	case 57:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:410
		{
			yyVAL.bytes = yyDollar[1].bytes
		}
	case 58:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:414
		{
			yyVAL.bytes = yyDollar[2].bytes
		}
	case 59:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:420
		{
			yyVAL.tableExprs = TableExprs{yyDollar[1].tableExpr}
		}
	case 60:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:424
		{
			yyVAL.tableExprs = append(yyVAL.tableExprs, yyDollar[3].tableExpr)
		}
	case 61:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:430
		{
			yyVAL.tableExpr = &AliasedTableExpr{Expr: yyDollar[1].smTableExpr, As: yyDollar[2].bytes, Hints: yyDollar[3].indexHints}
		}
	case 62:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:434
		{
			yyVAL.tableExpr = &ParenTableExpr{Expr: yyDollar[2].tableExpr}
		}
	case 63:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:438
		{
			yyVAL.tableExpr = &JoinTableExpr{LeftExpr: yyDollar[1].tableExpr, Join: yyDollar[2].str, RightExpr: yyDollar[3].tableExpr}
		}
	case 64:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line sql.y:442
		{
			yyVAL.tableExpr = &JoinTableExpr{LeftExpr: yyDollar[1].tableExpr, Join: yyDollar[2].str, RightExpr: yyDollar[3].tableExpr, On: yyDollar[5].boolExpr}
		}
	case 65:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:447
		{
			yyVAL.bytes = nil
		}
	case 66:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:451
		{
			yyVAL.bytes = yyDollar[1].bytes
		}
	case 67:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:455
		{
			yyVAL.bytes = yyDollar[2].bytes
		}
	case 68:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:461
		{
			yyVAL.str = AST_JOIN
		}
	case 69:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:465
		{
			yyVAL.str = AST_STRAIGHT_JOIN
		}
	case 70:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:469
		{
			yyVAL.str = AST_LEFT_JOIN
		}
	case 71:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:473
		{
			yyVAL.str = AST_LEFT_JOIN
		}
	case 72:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:477
		{
			yyVAL.str = AST_RIGHT_JOIN
		}
	case 73:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:481
		{
			yyVAL.str = AST_RIGHT_JOIN
		}
	case 74:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:485
		{
			yyVAL.str = AST_JOIN
		}
	case 75:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:489
		{
			yyVAL.str = AST_CROSS_JOIN
		}
	case 76:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:493
		{
			yyVAL.str = AST_NATURAL_JOIN
		}
	case 77:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:499
		{
			yyVAL.smTableExpr = &TableName{Name: yyDollar[1].bytes}
		}
	case 78:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:503
		{
			yyVAL.smTableExpr = &TableName{Qualifier: yyDollar[1].bytes, Name: yyDollar[3].bytes}
		}
	case 79:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:507
		{
			yyVAL.smTableExpr = yyDollar[1].subquery
		}
	case 80:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:513
		{
			yyVAL.tableName = &TableName{Name: yyDollar[1].bytes}
		}
	case 81:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:517
		{
			yyVAL.tableName = &TableName{Qualifier: yyDollar[1].bytes, Name: yyDollar[3].bytes}
		}
	case 82:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:522
		{
			yyVAL.indexHints = nil
		}
	case 83:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line sql.y:526
		{
			yyVAL.indexHints = &IndexHints{Type: AST_USE, Indexes: yyDollar[4].bytes2}
		}
	case 84:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line sql.y:530
		{
			yyVAL.indexHints = &IndexHints{Type: AST_IGNORE, Indexes: yyDollar[4].bytes2}
		}
	case 85:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line sql.y:534
		{
			yyVAL.indexHints = &IndexHints{Type: AST_FORCE, Indexes: yyDollar[4].bytes2}
		}
	case 86:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:540
		{
			yyVAL.bytes2 = [][]byte{yyDollar[1].bytes}
		}
	case 87:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:544
		{
			yyVAL.bytes2 = append(yyDollar[1].bytes2, yyDollar[3].bytes)
		}
	case 88:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:549
		{
			yyVAL.boolExpr = nil
		}
	case 89:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:553
		{
			yyVAL.boolExpr = yyDollar[2].boolExpr
		}
	case 91:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:560
		{
			yyVAL.boolExpr = &AndExpr{Left: yyDollar[1].boolExpr, Right: yyDollar[3].boolExpr}
		}
	case 92:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:564
		{
			yyVAL.boolExpr = &OrExpr{Left: yyDollar[1].boolExpr, Right: yyDollar[3].boolExpr}
		}
	case 93:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:568
		{
			yyVAL.boolExpr = &NotExpr{Expr: yyDollar[2].boolExpr}
		}
	case 94:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:572
		{
			yyVAL.boolExpr = &ParenBoolExpr{Expr: yyDollar[2].boolExpr}
		}
	case 95:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:578
		{
			yyVAL.boolExpr = &ComparisonExpr{Left: yyDollar[1].valExpr, Operator: yyDollar[2].str, Right: yyDollar[3].valExpr}
		}
	case 96:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:582
		{
			yyVAL.boolExpr = &ComparisonExpr{Left: yyDollar[1].valExpr, Operator: AST_IN, Right: yyDollar[3].colTuple}
		}
	case 97:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line sql.y:586
		{
			yyVAL.boolExpr = &ComparisonExpr{Left: yyDollar[1].valExpr, Operator: AST_NOT_IN, Right: yyDollar[4].colTuple}
		}
	case 98:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:590
		{
			yyVAL.boolExpr = &ComparisonExpr{Left: yyDollar[1].valExpr, Operator: AST_LIKE, Right: yyDollar[3].valExpr}
		}
	case 99:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line sql.y:594
		{
			yyVAL.boolExpr = &ComparisonExpr{Left: yyDollar[1].valExpr, Operator: AST_NOT_LIKE, Right: yyDollar[4].valExpr}
		}
	case 100:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line sql.y:598
		{
			yyVAL.boolExpr = &RangeCond{Left: yyDollar[1].valExpr, Operator: AST_BETWEEN, From: yyDollar[3].valExpr, To: yyDollar[5].valExpr}
		}
	case 101:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line sql.y:602
		{
			yyVAL.boolExpr = &RangeCond{Left: yyDollar[1].valExpr, Operator: AST_NOT_BETWEEN, From: yyDollar[4].valExpr, To: yyDollar[6].valExpr}
		}
	case 102:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:606
		{
			yyVAL.boolExpr = &NullCheck{Operator: AST_IS_NULL, Expr: yyDollar[1].valExpr}
		}
	case 103:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line sql.y:610
		{
			yyVAL.boolExpr = &NullCheck{Operator: AST_IS_NOT_NULL, Expr: yyDollar[1].valExpr}
		}
	case 104:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:614
		{
			yyVAL.boolExpr = &ExistsExpr{Subquery: yyDollar[2].subquery}
		}
	case 105:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line sql.y:618
		{
			yyVAL.boolExpr = &KeyrangeExpr{Start: yyDollar[3].valExpr, End: yyDollar[5].valExpr}
		}
	case 106:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:624
		{
			yyVAL.str = AST_EQ
		}
	case 107:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:628
		{
			yyVAL.str = AST_LT
		}
	case 108:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:632
		{
			yyVAL.str = AST_GT
		}
	case 109:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:636
		{
			yyVAL.str = AST_LE
		}
	case 110:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:640
		{
			yyVAL.str = AST_GE
		}
	case 111:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:644
		{
			yyVAL.str = AST_NE
		}
	case 112:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:648
		{
			yyVAL.str = AST_NSE
		}
	case 113:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:654
		{
			yyVAL.colTuple = ValTuple(yyDollar[2].valExprs)
		}
	case 114:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:658
		{
			yyVAL.colTuple = yyDollar[1].subquery
		}
	case 115:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:662
		{
			yyVAL.colTuple = ListArg(yyDollar[1].bytes)
		}
	case 116:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:668
		{
			yyVAL.subquery = &Subquery{yyDollar[2].selStmt}
		}
	case 117:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:674
		{
			yyVAL.valExprs = ValExprs{yyDollar[1].valExpr}
		}
	case 118:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:678
		{
			yyVAL.valExprs = append(yyDollar[1].valExprs, yyDollar[3].valExpr)
		}
	case 119:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:684
		{
			yyVAL.valExpr = yyDollar[1].valExpr
		}
	case 120:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:688
		{
			yyVAL.valExpr = yyDollar[1].colName
		}
	case 121:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:692
		{
			yyVAL.valExpr = yyDollar[1].rowTuple
		}
	case 122:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:696
		{
			yyVAL.valExpr = &BinaryExpr{Left: yyDollar[1].valExpr, Operator: AST_BITAND, Right: yyDollar[3].valExpr}
		}
	case 123:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:700
		{
			yyVAL.valExpr = &BinaryExpr{Left: yyDollar[1].valExpr, Operator: AST_BITOR, Right: yyDollar[3].valExpr}
		}
	case 124:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:704
		{
			yyVAL.valExpr = &BinaryExpr{Left: yyDollar[1].valExpr, Operator: AST_BITXOR, Right: yyDollar[3].valExpr}
		}
	case 125:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:708
		{
			yyVAL.valExpr = &BinaryExpr{Left: yyDollar[1].valExpr, Operator: AST_PLUS, Right: yyDollar[3].valExpr}
		}
	case 126:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:712
		{
			yyVAL.valExpr = &BinaryExpr{Left: yyDollar[1].valExpr, Operator: AST_MINUS, Right: yyDollar[3].valExpr}
		}
	case 127:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:716
		{
			yyVAL.valExpr = &BinaryExpr{Left: yyDollar[1].valExpr, Operator: AST_MULT, Right: yyDollar[3].valExpr}
		}
	case 128:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:720
		{
			yyVAL.valExpr = &BinaryExpr{Left: yyDollar[1].valExpr, Operator: AST_DIV, Right: yyDollar[3].valExpr}
		}
	case 129:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:724
		{
			yyVAL.valExpr = &BinaryExpr{Left: yyDollar[1].valExpr, Operator: AST_MOD, Right: yyDollar[3].valExpr}
		}
	case 130:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:728
		{
			if num, ok := yyDollar[2].valExpr.(NumVal); ok {
				switch yyDollar[1].byt {
				case '-':
					yyVAL.valExpr = append(NumVal("-"), num...)
				case '+':
					yyVAL.valExpr = num
				default:
					yyVAL.valExpr = &UnaryExpr{Operator: yyDollar[1].byt, Expr: yyDollar[2].valExpr}
				}
			} else {
				yyVAL.valExpr = &UnaryExpr{Operator: yyDollar[1].byt, Expr: yyDollar[2].valExpr}
			}
		}
	case 131:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:743
		{
			yyVAL.valExpr = &FuncExpr{Name: yyDollar[1].bytes}
		}
	case 132:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line sql.y:747
		{
			yyVAL.valExpr = &FuncExpr{Name: yyDollar[1].bytes, Exprs: yyDollar[3].selectExprs}
		}
	case 133:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line sql.y:751
		{
			yyVAL.valExpr = &FuncExpr{Name: yyDollar[1].bytes, Distinct: true, Exprs: yyDollar[4].selectExprs}
		}
	case 134:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line sql.y:755
		{
			yyVAL.valExpr = &FuncExpr{Name: yyDollar[1].bytes, Exprs: yyDollar[3].selectExprs}
		}
	case 135:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:759
		{
			yyVAL.valExpr = yyDollar[1].caseExpr
		}
	case 136:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:765
		{
			yyVAL.bytes = IF_BYTES
		}
	case 137:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:769
		{
			yyVAL.bytes = VALUES_BYTES
		}
	case 138:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:775
		{
			yyVAL.byt = AST_UPLUS
		}
	case 139:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:779
		{
			yyVAL.byt = AST_UMINUS
		}
	case 140:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:783
		{
			yyVAL.byt = AST_TILDA
		}
	case 141:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line sql.y:789
		{
			yyVAL.caseExpr = &CaseExpr{Expr: yyDollar[2].valExpr, Whens: yyDollar[3].whens, Else: yyDollar[4].valExpr}
		}
	case 142:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:794
		{
			yyVAL.valExpr = nil
		}
	case 143:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:798
		{
			yyVAL.valExpr = yyDollar[1].valExpr
		}
	case 144:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:804
		{
			yyVAL.whens = []*When{yyDollar[1].when}
		}
	case 145:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:808
		{
			yyVAL.whens = append(yyDollar[1].whens, yyDollar[2].when)
		}
	case 146:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line sql.y:814
		{
			yyVAL.when = &When{Cond: yyDollar[2].boolExpr, Val: yyDollar[4].valExpr}
		}
	case 147:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:819
		{
			yyVAL.valExpr = nil
		}
	case 148:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:823
		{
			yyVAL.valExpr = yyDollar[2].valExpr
		}
	case 149:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:829
		{
			yyVAL.colName = &ColName{Name: yyDollar[1].bytes}
		}
	case 150:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:833
		{
			yyVAL.colName = &ColName{Qualifier: yyDollar[1].bytes, Name: yyDollar[3].bytes}
		}
	case 151:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:839
		{
			yyVAL.valExpr = StrVal(yyDollar[1].bytes)
		}
	case 152:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:843
		{
			yyVAL.valExpr = NumVal(yyDollar[1].bytes)
		}
	case 153:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:847
		{
			yyVAL.valExpr = ValArg(yyDollar[1].bytes)
		}
	case 154:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:851
		{
			yyVAL.valExpr = &NullVal{}
		}
	case 155:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:856
		{
			yyVAL.valExprs = nil
		}
	case 156:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:860
		{
			yyVAL.valExprs = yyDollar[3].valExprs
		}
	case 157:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:865
		{
			yyVAL.boolExpr = nil
		}
	case 158:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:869
		{
			yyVAL.boolExpr = yyDollar[2].boolExpr
		}
	case 159:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:874
		{
			yyVAL.orderBy = nil
		}
	case 160:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:878
		{
			yyVAL.orderBy = yyDollar[3].orderBy
		}
	case 161:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:884
		{
			yyVAL.orderBy = OrderBy{yyDollar[1].order}
		}
	case 162:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:888
		{
			yyVAL.orderBy = append(yyDollar[1].orderBy, yyDollar[3].order)
		}
	case 163:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:894
		{
			yyVAL.order = &Order{Expr: yyDollar[1].valExpr, Direction: yyDollar[2].str}
		}
	case 164:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:899
		{
			yyVAL.str = AST_ASC
		}
	case 165:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:903
		{
			yyVAL.str = AST_ASC
		}
	case 166:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:907
		{
			yyVAL.str = AST_DESC
		}
	case 167:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:912
		{
			yyVAL.limit = nil
		}
	case 168:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:916
		{
			yyVAL.limit = &Limit{Rowcount: yyDollar[2].valExpr}
		}
	case 169:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line sql.y:920
		{
			yyVAL.limit = &Limit{Offset: yyDollar[2].valExpr, Rowcount: yyDollar[4].valExpr}
		}
	case 170:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:925
		{
			yyVAL.str = ""
		}
	case 171:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:929
		{
			yyVAL.str = AST_FOR_UPDATE
		}
	case 172:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line sql.y:933
		{
			if !bytes.Equal(yyDollar[3].bytes, SHARE) {
				yylex.Error("expecting share")
				return 1
			}
			if !bytes.Equal(yyDollar[4].bytes, MODE) {
				yylex.Error("expecting mode")
				return 1
			}
			yyVAL.str = AST_SHARE_MODE
		}
	case 173:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:946
		{
			yyVAL.columns = nil
		}
	case 174:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:950
		{
			yyVAL.columns = yyDollar[2].columns
		}
	case 175:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:956
		{
			yyVAL.columns = Columns{&NonStarExpr{Expr: yyDollar[1].colName}}
		}
	case 176:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:960
		{
			yyVAL.columns = append(yyVAL.columns, &NonStarExpr{Expr: yyDollar[3].colName})
		}
	case 177:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:965
		{
			yyVAL.updateExprs = nil
		}
	case 178:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line sql.y:969
		{
			yyVAL.updateExprs = yyDollar[5].updateExprs
		}
	case 179:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:975
		{
			yyVAL.insRows = yyDollar[2].values
		}
	case 180:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:979
		{
			yyVAL.insRows = yyDollar[1].selStmt
		}
	case 181:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:985
		{
			yyVAL.values = Values{yyDollar[1].rowTuple}
		}
	case 182:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:989
		{
			yyVAL.values = append(yyDollar[1].values, yyDollar[3].rowTuple)
		}
	case 183:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:995
		{
			yyVAL.rowTuple = ValTuple(yyDollar[2].valExprs)
		}
	case 184:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:999
		{
			yyVAL.rowTuple = yyDollar[1].subquery
		}
	case 185:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:1005
		{
			yyVAL.updateExprs = UpdateExprs{yyDollar[1].updateExpr}
		}
	case 186:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:1009
		{
			yyVAL.updateExprs = append(yyDollar[1].updateExprs, yyDollar[3].updateExpr)
		}
	case 187:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:1015
		{
			yyVAL.updateExpr = &UpdateExpr{Name: yyDollar[1].colName, Expr: yyDollar[3].valExpr}
		}
	case 188:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:1020
		{
			yyVAL.empty = struct{}{}
		}
	case 189:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:1022
		{
			yyVAL.empty = struct{}{}
		}
	case 190:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:1025
		{
			yyVAL.empty = struct{}{}
		}
	case 191:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line sql.y:1027
		{
			yyVAL.empty = struct{}{}
		}
	case 192:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:1030
		{
			yyVAL.empty = struct{}{}
		}
	case 193:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:1032
		{
			yyVAL.empty = struct{}{}
		}
	case 194:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:1036
		{
			yyVAL.empty = struct{}{}
		}
	case 195:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:1038
		{
			yyVAL.empty = struct{}{}
		}
	case 196:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:1040
		{
			yyVAL.empty = struct{}{}
		}
	case 197:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:1042
		{
			yyVAL.empty = struct{}{}
		}
	case 198:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:1044
		{
			yyVAL.empty = struct{}{}
		}
	case 199:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:1047
		{
			yyVAL.empty = struct{}{}
		}
	case 200:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:1049
		{
			yyVAL.empty = struct{}{}
		}
	case 201:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:1052
		{
			yyVAL.empty = struct{}{}
		}
	case 202:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:1054
		{
			yyVAL.empty = struct{}{}
		}
	case 203:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:1057
		{
			yyVAL.empty = struct{}{}
		}
	case 204:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line sql.y:1059
		{
			yyVAL.empty = struct{}{}
		}
	case 205:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line sql.y:1063
		{
			yyVAL.bytes = bytes.ToLower(yyDollar[1].bytes)
		}
	case 206:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line sql.y:1068
		{
			ForceEOF(yylex)
		}
	}
	goto yystack /* stack new state and value */
}
