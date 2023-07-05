pipeline {
    agent any
    
    tools { 
        nodejs "node"
        go "golang"
    }

    options {
        disableConcurrentBuilds()
    }
    
    stages {
        stage("Init") {
            steps {
                echo 'Installing dependencies...'
                dir("frontend") {
                    sh 'npm install'
                }
            }
        }
        stage("Build") {
            steps {
                echo 'Building...'
                dir("frontend") {
                    sh 'CI= npm run build'
                }
                dir("backend") {
                    sh 'go build -o dmrbac-exec'
                }
            }
        }
//         stage("test") {
//           steps {
//                 echo 'testing...'
//             }
//         }
        stage("Deploy") {
            steps {
                echo 'Starting service...'
//                 echo 'Deploying...'
                dir("backend") {
                    sh 'sudo systemctl restart dmrbac'
                }
            }
        }
    }   
}
