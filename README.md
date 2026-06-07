<div align="center">
  <img src="CerebroDigitalLogo.png" alt="Cerebro Digital Logo" width="400"/>
  <p><em>Plataforma de gestión de notas basada en grafos e IA local</em></p>
</div>

## Badges
![Estado del proyecto](https://img.shields.io/badge/Estado-En%20desarrollo-brightgreen)
![Version](https://img.shields.io/badge/Versión-1.0.0-blue)
![Issues abiertos](https://img.shields.io/github/issues/esf874/CerebroDigital)
![Último commit](https://img.shields.io/github/last-commit/esf874/CerebroDigital)
![Licencia](https://img.shields.io/badge/Licencia-MIT-green)

## Tecnologías
![Go](https://img.shields.io/badge/Go-00ADD8?style=flat-square&logo=go&logoColor=white)
![React](https://img.shields.io/badge/React-20232A?style=flat-square&logo=react&logoColor=61DAFB)
![JavaScript](https://img.shields.io/badge/JavaScript-F7DF1E?style=flat-square&logo=javascript&logoColor=black)
![MongoDB](https://img.shields.io/badge/MongoDB-47A248?style=flat-square&logo=mongodb&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-2496ED?style=flat-square&logo=docker&logoColor=white)
![Nix](https://img.shields.io/badge/Nix-5277C3?style=flat-square&logo=nixos&logoColor=white)

## Descripción
**Cerebro Digital** es una plataforma de gestión de conocimiento personal inspirada en el concepto de *Second Brain*. Permite almacenar, relacionar y consultar notas en formato markdown mediante una estructura basada en grafos.

La aplicación incorpora además un asistente de inteligencia artificial capaz de responder preguntas utilizando exclusivamente la información almacenada por el usuario mediante una arquitectura **Retrieval-Augmented Generation (RAG)** ejecutada íntegramente en local.

>El proyecto ha sido desarrollado como Trabajo Fin de Grado en Ingeniería Informática.

## Características principales

- 📝 Gestión de notas en formato Markdown.
- 🏷️ Organización mediante etiquetas y enlaces.
- 🕸️ Visualización interactiva del grafo de conocimiento.
- 🔗 Navegación contextual entre notas relacionadas.
- 🤖 Asistente conversacional basado en LLM.
- 🧠 Recuperación contextual mediante Graph-RAG.
- 🔒 Ejecución local y privada.
- 🏛️ Arquitectura basada en Clean Architecture.
- 🐳 Despliegue mediante Docker Compose.

## Arquitectura

### Componentes principales

| Componente | Tecnología | Descripción |
|------------|------------|-------------|
| **Frontend** | React | Interfaz de usuario interactiva |
| **API REST** | Go | Lógica de negocio y servicios |
| **Base de datos** | MongoDB | Persistencia de notas y relaciones |
| **Motor de inferencia** | llama.cpp | Ejecución local del modelo IA |
| **Modelo LLM** | Qwen3-4B-Thinking | Modelo de lenguaje para el asistente |
| **Sistema RAG** | Graph-RAG | Recuperación contextual basada en grafos |

### Diagrama de arquitectura

<img src="DiagramaArquitectura.png" alt="Diagrama de Arquitectura" width="700" />

## Tecnologías utilizadas

### Backend
- **Go** - Lenguaje principal del servidor
- **MongoDB** - Base de datos NoSQL
- **Docker** - Contenerización

### Frontend
- **React** - Biblioteca para la interfaz de usuario
- **HTML5** - Estructura del documento
- **CSS3** - Estilos y diseño responsive
- **JavaScript** - Lógica del cliente

### Inteligencia Artificial
- **llama.cpp** - Inferencia eficiente de LLMs en CPU
- **Qwen3-4B-Thinking** - Modelo de lenguaje optimizado
- **Graph-RAG** - Recuperación aumentada por grafos

### DevOps y calidad
- **Docker Compose** - Orquestación de servicios
- **SonarCloud** - Análisis estático de código
- **GitLab** - Control de versiones y CI/CD

## Capturas de pantalla

### Vista principal

<img src="PantallaPrincipal.png" alt="Vista principal" width="600" />

### Vista detallada de una nota

<img src="VistaDetallada.png" alt="Vista detallada de una nota" width="600" />

### Visualización del grafo

<img src="Subgrafo.png" alt="Subgrafo de una nota concreta" width="600" />

### Asistente conversacional

<img src="AsistenteGlobal.png" alt="Asistente global" width="600" />


## Instalación 

### Requisitos previos
- [Docker](https://docs.docker.com/get-docker/) (versión 20.10+)
- [Docker Compose](https://docs.docker.com/compose/install/) (versión 2.0+)
- [Git](https://git-scm.com/downloads)

1. **Clonar el repositorio**
    ```bash
   git clone https://github.com/esf874/CerebroDigital
   cd CerebroDigital
   
2. **Iniciar la aplicación**
    docker compose up --build

3. **Acceder a los servicios**
    Una vez desplegado:
    
        Frontend: http://localhost:8080 
        Backend: http://localhost:8080 
        MongoDB: localhost:27017

## Configuración del modelo IA
El sistema utiliza inferencia local mediante **llama.cpp** y el modelo **Qwen3-4B-Thinking**.

### Requisitos previos
 1. Descargar el modelo en formato GGUF
 2. Ubicarlo en el directorio configurado para llama.cpp

## Calidad del software

La calidad del código ha sido evaluada mediante **SonarCloud**, herramienta de análisis estático ampliamente utilizada en entornos profesionales para la detección de defectos, vulnerabilidades y problemas de mantenibilidad.

### Resultados del análisis
El análisis final del proyecto supera satisfactoriamente el Quality Gate establecido, obteniendo los siguientes resultados:

| Métrica | Estado |
|---------|--------|
| ✅ **Quality Gate** | **Passed** |
| 🔒 **Security** | **A** |
| 🔧 **Maintainability** | **A** |
| 🐛 **Reliability** | **B** |
| ⚠️ **Vulnerabilities** | **0** |
| 📝 **Duplicated Code** | **2.2%** |

<img src="ResumenSonar.png" alt="Resumen Análisis SonarCloud" width="500" />

## Autor
**Estela Simón Fernández**
