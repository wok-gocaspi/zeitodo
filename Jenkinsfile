pipeline {
    agent any

    tools {
        go 'go-1.18.5'
    }
    environment {
        GO111MODULE = 'on'
        CGO_ENABLED = 0
    }

    parameters {
            gitParameter branchFilter: 'origin.*/(.*)',
                         name: 'BRANCH_TAG',
                         type: 'PT_BRANCH_TAG',
                         defaultValue: 'main'
        }

    stages {
        stage('Build') {
            steps {
                echo 'Building..'
                sh " go version"
                checkout([$class: 'GitSCM',
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
        stage('Docker Image') {
            steps {
                echo 'Building Docker Image....'
                script{
                    docker.build("zeitodoreg.azurecr.io/zeitodobackend:${env.BRANCH_TAG}-${env.BUILD_NUMBER}")
                }
            }
        }
    }
}
