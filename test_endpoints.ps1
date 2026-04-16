# ================================================================
# COMPREHENSIVE ENDPOINT TEST v2 - After Fixes
# ================================================================
$ErrorActionPreference = "Continue"

$base = "http://localhost:8080/api/v1"
$results = @()

function Test-Endpoint {
    param(
        [string]$Name,
        [string]$Method,
        [string]$Uri,
        [hashtable]$Headers = @{},
        [string]$Body = "",
        [int]$ExpectedStatus = 200
    )
    
    try {
        $params = @{ Uri = $Uri; Method = $Method; Headers = $Headers; UseBasicParsing = $true }
        if ($Body) {
            $params["ContentType"] = "application/json"
            $params["Body"] = $Body
        }
        
        $res = Invoke-WebRequest @params
        $json = $res.Content | ConvertFrom-Json
        $status = $json.status
        $msg = $json.message
        
        if ($status -eq $ExpectedStatus) {
            Write-Host "[PASS] $Name" -ForegroundColor Green
            Write-Host "       $msg" -ForegroundColor DarkGray
            return @{ Name=$Name; Status="PASS"; Code=$status; Message=$msg; Data=$json.data }
        } else {
            Write-Host "[WARN] $Name - status $status (expected $ExpectedStatus)" -ForegroundColor Yellow
            Write-Host "       $msg" -ForegroundColor DarkGray
            return @{ Name=$Name; Status="WARN"; Code=$status; Message=$msg; Data=$json.data }
        }
    } catch {
        $errBody = $_.ErrorDetails.Message
        $httpStatus = 0; $errMsg = ""
        if ($errBody) {
            try { $errJson = $errBody | ConvertFrom-Json; $httpStatus = $errJson.status; $errMsg = $errJson.message } catch { $errMsg = $errBody }
        } else { $errMsg = $_.Exception.Message }
        
        Write-Host "[FAIL] $Name" -ForegroundColor Red
        Write-Host "       HTTP $httpStatus | $errMsg" -ForegroundColor DarkGray
        if ($errBody -and $errBody.Length -lt 300) { Write-Host "       $errBody" -ForegroundColor DarkRed }
        return @{ Name=$Name; Status="FAIL"; Code=$httpStatus; Message=$errMsg; RawError=$errBody }
    }
}

# ================================================================
# LOGIN
# ================================================================
Write-Host "`n==================== AUTH ====================" -ForegroundColor Magenta
$loginRes = Test-Endpoint -Name "Login" -Method "POST" -Uri "$base/users/login" `
    -Body '{"email":"testuser@example.com","password":"Password123"}'
$token = $loginRes.Data
$auth = @{ "Authorization" = "Bearer $token" }
$results += $loginRes

# ================================================================
# USERS
# ================================================================
Write-Host "`n==================== USERS ====================" -ForegroundColor Magenta
$results += Test-Endpoint -Name "Users: GetProfile" -Method "GET" -Uri "$base/users/profile" -Headers $auth
$results += Test-Endpoint -Name "Users: GetAllUser" -Method "GET" -Uri "$base/users/all" -Headers $auth

# ================================================================
# CATEGORIES
# ================================================================
Write-Host "`n==================== CATEGORIES ====================" -ForegroundColor Magenta
$results += Test-Endpoint -Name "Categories: Create (income)" -Method "POST" -Uri "$base/categories" -Headers $auth `
    -Body '{"name":"Gaji Test","type":"income"}'
$results += Test-Endpoint -Name "Categories: Create (expense)" -Method "POST" -Uri "$base/categories" -Headers $auth `
    -Body '{"name":"Makan Test","type":"expense"}'
$results += Test-Endpoint -Name "Categories: GetAll" -Method "GET" -Uri "$base/categories/all?page=1&limit=10" -Headers $auth
$results += Test-Endpoint -Name "Categories: Update" -Method "PUT" -Uri "$base/categories/update" -Headers $auth `
    -Body '{"name":"Gaji Updated"}'

# Get category IDs
try {
    $catRes = Invoke-RestMethod -Uri "$base/categories/all?page=1&limit=10" -Method GET -Headers $auth
    $catIncome = ($catRes.data.items | Where-Object { $_.type -eq "income" } | Select-Object -First 1).id
    $catExpense = ($catRes.data.items | Where-Object { $_.type -eq "expense" } | Select-Object -First 1).id
    if (-not $catIncome) { $catIncome = $catRes.data.items[0].id }
    if (-not $catExpense) { $catExpense = $catRes.data.items[0].id }
    Write-Host "       catIncome=$catIncome | catExpense=$catExpense" -ForegroundColor DarkGray
} catch { $catIncome = $null; $catExpense = $null }

if ($catIncome) {
    $results += Test-Endpoint -Name "Categories: GetById" -Method "GET" -Uri "$base/categories/$catIncome" -Headers $auth
}

# ================================================================
# TRANSACTIONS
# ================================================================
Write-Host "`n==================== TRANSACTIONS ====================" -ForegroundColor Magenta

# Create income
if ($catIncome) {
    $txBody = @{ type="income"; amount=5000000; category_id=$catIncome; description="Gaji april"; date="2026-04-15T00:00:00Z" } | ConvertTo-Json
    $results += Test-Endpoint -Name "Tx: Create (income)" -Method "POST" -Uri "$base/transactions" -Headers $auth -Body $txBody
}

# Create expense
if ($catExpense) {
    $txBodyExp = @{ type="expense"; amount=100000; category_id=$catExpense; description="Makan siang"; date="2026-04-15T00:00:00Z" } | ConvertTo-Json
    $results += Test-Endpoint -Name "Tx: Create (expense)" -Method "POST" -Uri "$base/transactions" -Headers $auth -Body $txBodyExp
}

$results += Test-Endpoint -Name "Tx: GetAll" -Method "GET" -Uri "$base/transactions/all?page=1&limit=10" -Headers $auth
$results += Test-Endpoint -Name "Tx: Update" -Method "PUT" -Uri "$base/transactions/update" -Headers $auth -Body '{"amount":6000000}'

# Analytics
$results += Test-Endpoint -Name "Tx: AvgIncomeDay" -Method "GET" -Uri "$base/transactions/avg-income-day" -Headers $auth
$results += Test-Endpoint -Name "Tx: AvgExpenseDay" -Method "GET" -Uri "$base/transactions/avg-expense-day" -Headers $auth
$results += Test-Endpoint -Name "Tx: AvgIncomeWeek" -Method "GET" -Uri "$base/transactions/avg-income-week" -Headers $auth
$results += Test-Endpoint -Name "Tx: AvgExpenseWeek" -Method "GET" -Uri "$base/transactions/avg-expense-week" -Headers $auth
$results += Test-Endpoint -Name "Tx: AvgIncomeMonth" -Method "GET" -Uri "$base/transactions/avg-income-month" -Headers $auth
$results += Test-Endpoint -Name "Tx: AvgExpenseMonth" -Method "GET" -Uri "$base/transactions/avg-expense-month" -Headers $auth

# Expense/Income by type (now uses JSON body)
$results += Test-Endpoint -Name "Tx: ExpenseByType" -Method "GET" -Uri "$base/transactions/expense" -Headers $auth -Body '{"type":"expense"}'
$results += Test-Endpoint -Name "Tx: IncomeByType" -Method "GET" -Uri "$base/transactions/income" -Headers $auth -Body '{"type":"income"}'

$results += Test-Endpoint -Name "Tx: ExpenseDayCategory" -Method "GET" -Uri "$base/transactions/expense-day-category" -Headers $auth
$results += Test-Endpoint -Name "Tx: IncomeDayCategory" -Method "GET" -Uri "$base/transactions/income-day-category" -Headers $auth
$results += Test-Endpoint -Name "Tx: TotalExpenseDay" -Method "GET" -Uri "$base/transactions/total-expense-day" -Headers $auth
$results += Test-Endpoint -Name "Tx: TotalExpenseWeek" -Method "GET" -Uri "$base/transactions/total-expense-week" -Headers $auth
$results += Test-Endpoint -Name "Tx: TotalExpenseMonth" -Method "GET" -Uri "$base/transactions/total-expense-month" -Headers $auth
$results += Test-Endpoint -Name "Tx: TotalIncomeDay" -Method "GET" -Uri "$base/transactions/total-income-day" -Headers $auth
$results += Test-Endpoint -Name "Tx: TotalIncomeWeek" -Method "GET" -Uri "$base/transactions/total-income-week" -Headers $auth
$results += Test-Endpoint -Name "Tx: TotalIncomeMonth" -Method "GET" -Uri "$base/transactions/total-income-month" -Headers $auth

# ================================================================
# BUDGETS
# ================================================================
Write-Host "`n==================== BUDGETS ====================" -ForegroundColor Magenta

if ($catIncome) {
    $budgetBody = @{
        category_id=$catIncome; limit_amount=1000000; period="monthly"
        start_date="2026-04-01T00:00:00Z"; end_date="2026-04-30T00:00:00Z"; is_active=$true
    } | ConvertTo-Json
    $results += Test-Endpoint -Name "Budgets: Create" -Method "POST" -Uri "$base/budgets" -Headers $auth -Body $budgetBody
}

$results += Test-Endpoint -Name "Budgets: GetAll" -Method "GET" -Uri "$base/budgets?page=1&limit=10" -Headers $auth
$results += Test-Endpoint -Name "Budgets: Update" -Method "PUT" -Uri "$base/budgets/update" -Headers $auth -Body '{"limit_amount":2000000}'

# ================================================================
# GOALS
# ================================================================
Write-Host "`n==================== GOALS ====================" -ForegroundColor Magenta

$goalsBody = @{
    name="Tabungan Rumah"; target_amount=100000000; current_amount=5000000
    start_date="2026-04-01"; target_date="2027-04-01"
} | ConvertTo-Json
$results += Test-Endpoint -Name "Goals: Create" -Method "POST" -Uri "$base/goals" -Headers $auth -Body $goalsBody

$results += Test-Endpoint -Name "Goals: GetAll" -Method "GET" -Uri "$base/goals?page=1&limit=10" -Headers $auth

$results += Test-Endpoint -Name "Goals: Update" -Method "PUT" -Uri "$base/goals/update" -Headers $auth `
    -Body '{"name":"Tabungan Updated","target_amount":200000000,"current_amount":10000000}'

$results += Test-Endpoint -Name "Goals: Progress" -Method "GET" -Uri "$base/goals/progress" -Headers $auth
$results += Test-Endpoint -Name "Goals: RemainingDays" -Method "GET" -Uri "$base/goals/remaining-days" -Headers $auth
$results += Test-Endpoint -Name "Goals: Delete" -Method "DELETE" -Uri "$base/goals/delete" -Headers $auth

# ================================================================
# SUMMARY
# ================================================================
Write-Host "`n================================================================" -ForegroundColor White
Write-Host "                    TEST RESULTS SUMMARY v2" -ForegroundColor White
Write-Host "================================================================" -ForegroundColor White

$pass = ($results | Where-Object { $_.Status -eq "PASS" }).Count
$fail = ($results | Where-Object { $_.Status -eq "FAIL" }).Count
$warn = ($results | Where-Object { $_.Status -eq "WARN" }).Count
$total = $results.Count

Write-Host "`nTotal: $total | PASS: $pass | FAIL: $fail | WARN: $warn" -ForegroundColor White

if ($fail -gt 0 -or $warn -gt 0) {
    Write-Host "`n--- FAILED/WARNING ENDPOINTS ---" -ForegroundColor Red
    $results | Where-Object { $_.Status -ne "PASS" } | ForEach-Object {
        $color = if ($_.Status -eq "FAIL") { "Red" } else { "Yellow" }
        Write-Host "[$($_.Status)] $($_.Name)" -ForegroundColor $color
        Write-Host "   HTTP $($_.Code): $($_.Message)" -ForegroundColor DarkGray
        if ($_.RawError) { Write-Host "   $($_.RawError)" -ForegroundColor DarkRed }
    }
}

Write-Host "`n--- PASSED ENDPOINTS ---" -ForegroundColor Green
$results | Where-Object { $_.Status -eq "PASS" } | ForEach-Object {
    Write-Host "[PASS] $($_.Name)" -ForegroundColor Green
}

Write-Host "`n================================================================" -ForegroundColor White
