pipeline{
    agent any

    options { buildDiscarder(logRotator(numToKeepStr: '5')) }

    environment {
        ROLLBAR_API_KEY = credentials("RespondologyRollbarReadOnlyApiKey")
    }
    stages{
        stage("Remove artificat"){
            steps{
                echo "================= Removing Go binary ================="
                sh "rm -f main"
            }
        }
        stage("Build binary"){
            steps{
                script{
                    echo "================= Building binary ================="
                    sh "go build main.go"
                }
            }
        }
        stage("Checking Rollbar API"){
            steps{
                script{
                    sh "./main ${ROLLBAR_API_KEY}"
                }
            }
        }
    }

    post{
        always{
            echo "======== Pipeline complete ========"
        }
        success{
            echo "======== pipeline executed successfully ========"

        }
        failure{
            echo "======== pipeline execution failed ========"
            slackSend (color: "danger", message: "FAILURE: Job '${env.JOB_NAME} [${env.BUILD_NUMBER}]' (${env.BUILD_URL}")
        }
    }
}