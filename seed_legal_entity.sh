#!/usr/bin/env bash
set -euo pipefail

CORE_URL="${CORE_URL:-http://localhost:8080/api/v1}"
AUTH_URL="${AUTH_URL:-http://localhost:8081/auth/v1}"
PROGRAM_ID="${PROGRAM_ID:-6834f2c1-47cf-478e-a0df-26017307c84d}"

USERNAME="${USERNAME:-1}"
PASSWORD="${PASSWORD:-1}"

LISTENERS_COUNT="${LISTENERS_COUNT:-2}"

need() { command -v "$1" >/dev/null 2>&1 || { echo "Missing dependency: $1"; exit 1; }; }
need curl
need jq
need uuidgen

urlencode() {
    printf '%s' "$1" | jq -sRr @uri
}

api_post() {
  local url="$1"
  local body="$2"
  curl -sS -X POST "$url" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d "$body"
}

api_get() {
  local url="$1"
  curl -sS "$url" -H "Authorization: Bearer $TOKEN"
}

RUN_ID="$(date -u +%Y%m%d-%H%M%S)-$RANDOM"
echo "RUN_ID=$RUN_ID"

echo "[1/5] Login..."
LOGIN_RESP="$(curl -sS -X POST "$AUTH_URL/login" \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$USERNAME\",\"password\":\"$PASSWORD\"}")"

TOKEN="$(echo "$LOGIN_RESP" | jq -r '.token.access_token // empty')"
if [[ -z "$TOKEN" || "$TOKEN" == "null" ]]; then
  echo "ERROR: couldn't extract access token from:"
  echo "$LOGIN_RESP" | jq .
  exit 1
fi
echo "OK: token acquired"

echo "[2/5] Create executor..."
EXECUTOR_PAYLOAD="$(jq -n --arg run "$RUN_ID" '{
  status: "Директор",
  first_name: "Иван",
  second_name: "Иванов",
  middle_name: "Иванович",
  doverenost: "-"
}')"

api_post "$CORE_URL/executer/" "$EXECUTOR_PAYLOAD" | jq -r '.message // .'
echo "OK: executor created"

echo "[3/5] Create legal entity..."

LE_NAME="ООО ТестЮрик ${RUN_ID}"

LE_PAYLOAD="$(jq -n --arg name "$LE_NAME" --arg run "$RUN_ID" '{
  legal_entity: {
    name_company: $name,
    inn: "7701234567",
    kpp: "770101001",
    ogrn: "1027700132195",
    phone: "+79990000000",
    email: ("legal_" + $run + "@test.local"),
    first_name: "Пётр",
    second_name: "Петров",
    middle_name: "Тестович",
    status: "Директор"
  },
  reg_address: {
    mail_index: "123456",
    region: "Москва",
    city: "Москва",
    street: "Тверская",
    house: "1",
    building: "1",
    apartment: "1"
  }
}')"

api_post "$CORE_URL/legalentity/" "$LE_PAYLOAD" >/dev/null || true

echo "Finding legal entity ID..."
LE_FILTER="$(urlencode "$LE_NAME")"
LE_LIST="$(api_get "$CORE_URL/legalentity/?page=1&filter=$LE_FILTER")"

LE_ID="$(echo "$LE_LIST" | jq -r --arg name "$LE_NAME" '
  .data[]?
  | select((.name_company? // .company_name? // "") == $name)
  | (.id_legalentity? // .idLegalEntity? // .id? // .uuid? // empty)
' | head -n 1)"

if [[ -z "$LE_ID" || "$LE_ID" == "null" ]]; then
  echo "ERROR: couldn't find legal entity id in response."
  echo "$LE_LIST" | jq .
  echo
  echo "Fix: tell me what fields legalentity GET returns (paste one item from .data[0])."
  exit 1
fi
echo "OK: LegalEntity ID = $LE_ID"

echo "[4/5] Create listeners..."

LISTENER_IDS=()

for i in $(seq 1 "$LISTENERS_COUNT"); do
  SNILS="$(printf "123-456-78%1d 0%1d" "$i" "$i")"
  PHONE="+79990000000"
  EMAIL="listener_${RUN_ID}_${i}@test.local"

  LISTENER_PAYLOAD="$(jq -n \
    --arg le "$LE_ID" \
    --arg fn "Слушатель${i}" \
    --arg sn "Тестовый" \
    --arg mn "$RUN_ID" \
    --arg dob "2000-01-0$i" \
    --arg snils "$SNILS" \
    --arg phone "$PHONE" \
    --arg email "$EMAIL" \
    '{
      listener: {
        first_name: $fn,
        second_name: $sn,
        middle_name: $mn,
        date_of_birth: $dob,
        snils: $snils,
        contact_phone: $phone,
        email: $email,
        id_legalentity: $le,
        id_contractor: "00000000-0000-0000-0000-000000000000",
        looting_education: false
      },
      registration_address: {
        mail_index: "123456",
        region: "Москва",
        city: "Москва",
        street: "Тверская",
        house: "1",
        building: "1",
        apartment: "1"
      },
      passport: {
        place_birth: "Москва",
        citizenship: "РФ",
        gender: "Муж",
        seria: "1234",
        number: "567890",
        passport_given: "ОВД",
        date_given: "2018-01-01",
        code: "123-456"
      }
    }')"

  RESP="$(api_post "$CORE_URL/listener/" "$LISTENER_PAYLOAD")"
  echo "Listener create response: $RESP"  echo "Listener $i created, finding ID..."

  L_FILTER="$(urlencode "$EMAIL")"
  L_LIST="$(api_get "$CORE_URL/listener/?page=1&filter=$L_FILTER")"

  LID="$(echo "$L_LIST" | jq -r --arg email "$EMAIL" '
    .data[]?
    | select((.email? // "") == $email)
    | (.id_listener? // .idListener? // .id? // .uuid? // empty)
  ' | head -n 1)"

  if [[ -z "$LID" || "$LID" == "null" ]]; then
    echo "ERROR: couldn't find listener id for $EMAIL"
    echo "$L_LIST" | jq .
    exit 1
  fi

  echo "OK: Listener $i ID = $LID"
  LISTENER_IDS+=("$LID")
done

echo "[5/5] Enroll listeners to program $PROGRAM_ID ..."
for LID in "${LISTENER_IDS[@]}"; do
  ENROLL_PAYLOAD="$(jq -n --arg l "$LID" --arg p "$PROGRAM_ID" '{
    id_listener: $l,
    id_program: $p,
    start_date: "2026-03-02",
    end_date: "2026-03-30",
    current_price: 10000,
    is_active: true,
    group: "E2E-A",
    type_of_retraining: "Повышение квалификации"
  }')"

  api_post "$CORE_URL/enrollment/" "$ENROLL_PAYLOAD" | jq -r '.message // .'
done

echo
echo "DONE ✅"
echo "LegalEntity ID: $LE_ID"
echo "Listener IDs: ${LISTENER_IDS[*]}"
echo "Program ID: $PROGRAM_ID"
