package benchmark

const RETRY = 2
const FORWARD_BROWSES = 10
const BACKWARD_BROWSE_INTERVAL = 3

const OPEN_ORDER_DEFAULT_CANCEL_GROUP_SIZE = 1
const OPEN_ORDER_BIG_CANCEL_GROUP_SIZE = 5
const VEHICLES_NOT_SOLD = 10

const MIN_ORDERLINE_COUNT = 1
const MAX_ORDERLINE_COUNT = 5
const REMAINING_ORDERLINE_COUNT = 1
const MIN_ORDER_VEHICLE_COUNT = 10
const MAX_ORDER_VEHICLE_COUNT = 20
const MIN_LARGE_ORDER_VEHICLE_COUNT = 105
const MAX_LARGE_ORDER_VEHICLE_COUNT = 200

const ORDERS_PER_TXRATE = 750
const CUSTOMERS_PER_SCALE = 7500
const ITEMS_PER_CATEGORY = 200
const PART_POSTFIX_LENGTH = 10

const DEFERRED_ORDER_PERCENTAGE = 0.5
const CHECKOUT_ENTIRE_CART_PERCENTAGE = 0.5
const REFILL_CART_PERCENTAGE = 0.8
const LARGE_ORDER_PERCENTAGE = 0.1
const PERCENTAGE_1 = 1.0

const PART_PREFIX = "00001MITEM"

var txRate int
var dbSize int
var maxItemsPerLoc int
var itemsPerTxRate int
var parallelism int
var audit bool
var stopIfAuditFailed bool
var plannedLineBorrowPercent int
var numOfItems int
var numItemsPerLoc int

const OperationBrowseId = 0
const OperationManageId = 1
const OperationPurchaseId = 2
