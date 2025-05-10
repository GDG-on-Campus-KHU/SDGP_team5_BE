<h1 align="center">🚨 ResQ 🚨</h1>

<pre><code>여행 중 <b>위급 상황</b> 발생 시 <b>신속한 대응</b>을 지원하는 모바일 애플리케이션입니다.
</code></pre>


## 🎯 Goal

<pre><code>예기치 못한 <b>응급 상황</b>은 언제, 어디서든 발생할 수 있습니다.<br>

<b>ResQ</b>는 그런 순간에 여러분의 <b>생명</b>과 <b>안전</b>을 지켜주는 <b>든든한 동반자</b>가 되고자 합니다.</code></pre>

<img src="https://drive.google.com/uc?id=1wOQnC2JamqkYSGTIdRcSwDGFH92ATsjX" width=60%>

---

### 💡 Key Feature
<h4> 🆘 위급 상황 시 신속한 도움 요청</h4>

<pre><code>* <b>👆한 번의 터치</b>로 <b>응급 전화</b> 연결 및 <b>상황 녹음</b> 시작
* 주요 응급 상황에 대한 <b>대처법</b> 안내</code></pre>

---

#### 📌 기본 기능

- **Google** 계정으로 로그인

- **응급 상황별 대처법** 안내 -- **`🌐ko(한국어)`**, **`🌐en(영어)`**

- 자주 사용하는 대처법 **'즐겨찾기'** 추가

- 메인 화면에 표시되지 않는 대처법은 **키워드로 검색** 가능

---

#### 📌 의료 정보 관리

- **알레르기**, **복용 중인 약** 등 응급 상황에 필요한 **의료 정보** 등록

- **여행 국가** 설정 -- 🌐 `KR`, `US`, `GB`, `JP`, `CN`, `DE`, `FR`, `MX`

- 함께 여행할 가족 및 친구를 **그룹**에 추가

- **일행**의 의료 정보를 **여행 국가의 언어**로 **번역**

- 신장, 체중은 설정한 **여행 국가**에 따라 단위 변환

---

#### 📌 상황 녹음 저장 및 관리

- 응급 상황에 자동 녹음된 **음성 파일** 저장

- 상황 녹음에 대한 **텍스트**로 저장 (`Speech-to-Text`)

---

<br>

## 👥 Team Members
📌 2025 파트 연합 장기 프로젝트 SDGP - team5

| Name       | Role     | GitHub                                               |
|------------|----------|------------------------------------------------------|
| 권동현      | Mobile  | [GwonDongHyeon21](https://github.com/GwonDongHyeon21) |
| 김민        | Backend | [kmin1231](https://github.com/kmin1231) |
| 김태훈      | Mobile  | [taeh-kim](https://github.com/taeh-kim) |
| 박상영      | Backend | [Imsyp](https://github.com/Imsyp) |

<br>

## 🧩 Project Architecture

<img src="https://drive.google.com/uc?id=1WH0xObPY-U4_opNcNv_3qroycj5ra3ev" width=80%>

<br>

## 📂 Project Structure

```
.
├── 🔒 auth/
├── ⭐ favorite/
├── 👥 group/
├── 🌍 language/
├── 💊 medical_info/
├── 🎤 recording/
├── 🚑 situation/
├── 👤 user/
├── 🔧 util/
├── db
│   ├── gcs.go
│   └── mongo.go
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── 🐋 Dockerfile
├── go.mod
├── go.sum
├── ▶️ main.go
└── README.md
```

<br>

## ▶️ How to Run

```
git clone https://github.com/GDG-on-Campus-KHU/SDGP_team5_BE.git
```

```
cd SDGP_team5_BE
```

```
go mod tidy
```

```
air
```

<br>
