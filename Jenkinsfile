pipeline {
    agent any

    tools {
        go 'go-1.16.2'
    }
    environment {
        GO111MODULE = 'on'
        CGO_ENABLED = 0
    }

    parameters {
            gitParameter name: 'BRANCH_TAG',
                         type: 'PT_BRANCH_TAG',
                         defaultValue: 'master'
        }

    stages {
        stage('Build') {
            steps {
                echo 'Building..'
                sh " go version"
                checkout([$class: 'GitSCM',
                    //branches: [[name: '*/main']],
                    branches: [[name: "${params.BRANCH_TAG}"]],
                    userRemoteConfigs: [[url: 'https://github.com/wok-gocaspi/zeitodo.git']]])
                sh "go build"
            }
        }
        stage('Test') {
            steps {
                echo 'Testing..'
                sh "go generate ./..."
                //sh "go get gotest.tools/gotestsum"
                //sh "go test -v 2>&1 ./... | go-junit-report -set-exit-code > test-report.xml"
                //junit testResults: 'test-report.xml'
                sh "go test ./..."
            }
        }
        stage('Deploy') {
            steps {
                echo 'Deploying....'
                script{
                    docker.build "zeitodo:${env.BUILD_TAG}"
                }
            }
        }
    }
}