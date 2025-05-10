<h1 align="center">🚨 ResQ 🚨</h1>

<div align="center">
  <pre><code><b><i>Instant support for your travel emergencies!</i></b></code></pre>
  <h4>🌐 This is the <code>English</code> version of the README. | <a href="README.ko.md">한국어 버전</a></h4>
</div>

## 🎯 Goal

<pre><code><b>Unexpected emergencies</b> can happen anywhere and at any time.<br>

<b>ResQ</b> aims to be a <b>reliable companion</b> that protects your <b>life</b> and <b>safety</b> in such moments.</code></pre>

<img src="https://drive.google.com/uc?id=1wOQnC2JamqkYSGTIdRcSwDGFH92ATsjX" width=60%>

---

### 💡 Key Features
<h4> 🆘 Quickly request help in emergency situations</h4>

<pre><code>* <b>👆One tap</b> to connect to an <b>emergency hotline</b> and start <b>recording</b> the situation
* Provides <b>guidelines</b> on how to handle common <b>emergency situations</b></code></pre>

---


#### 📌 Basic Features

- Log in with **Google** account

- **Emergency response guidelines** for each situation -- **`🌐ko(Korean)`**, **`🌐en(English)`**

- Add **favorite** emergency responses

- **Search** for emergency responses that are not displayed on the main screen

---

#### 📌 Medical Information Management

- Fill out **medical information** such as **allergies**, **medications**, etc. for emergencies

- Set **destination country** -- 🌐 `KR`, `US`, `GB`, `JP`, `CN`, `DE`, `FR`, `MX`

- Add family and friends who are traveling with you to the **group**

- **Translate medical information** of the group members into the **language** of the destination country

- Convert **height** and **weight** units based on the **destination country's** measurement system

---

#### 📌 Situation Recording

- Automatically save **audio files** recorded during emergencies

- Save situation recordings as **text** using the **`Speech-to-Text`** API

---

<br>

## 🧩 Project Architecture

<img src="https://drive.google.com/uc?id=1WH0xObPY-U4_opNcNv_3qroycj5ra3ev" width=80%>

<br>

## ⚙️ Tech Stacks

| stack                                     | description                                                     |
| ----------------------------------------- | --------------------------------------------------------------- |
| **`Go`**                                  | **main server** programming language                                  |                   
| `Gin`                                     | web framework for Go-based server                                 |
| **`MongoDB`**                             | **core database** for storing and managing application data          |
| **`GCP`** (Google Cloud Platform)         | **cloud platform** for operating the entire system infrastructure and services |
| **`Gemini`**                              | provides **medical context-based translations** of medical information |
| `Google Cloud Speech-to-Text API`         | converts recorded emergency situation audio files to **text**         |
| `Elasticsearch + Monstache`               | synchronizes MongoDB data to enable **fast search** |
| `Kibana`                                  | visualizes **Elasticsearch** data, provides data analysis and dashboards |
| `Docker` + `GitHub Actions`               | builds and deploys the server automatically via a **CI/CD** environment |
| `GAR` (Google Artifact Registry)          | stores Docker images for server deployment                       |
| `GCS` (Google Cloud Storage)              | temporarily stores automatically recorded audio files during emergencies |
| `Google Cloud Load Balancer`              | **distributes traffic** to ensure service load balancing and high availability |
| `Google Cloud DNS`                        | configures domain names to connect users to the service          |
| `Google Cloud CDN`                        | caches static content to reduce latency and improve response times |

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


## 👥 Team Members

| Name       | English Name   | Role     | GitHub                                                |
|------------|----------------|----------|-------------------------------------------------------|
| 권동현      | DongHyeon Gwon | Mobile   | [GwonDongHyeon21](https://github.com/GwonDongHyeon21) |
| 김민        | Min Kim        | Backend | [kmin1231](https://github.com/kmin1231) |
| 김태훈      | Taehoon Kim    | Mobile  | [taeh-kim](https://github.com/taeh-kim) |
| 박상영      | SangYeong Park | Backend | [Imsyp](https://github.com/Imsyp) |
