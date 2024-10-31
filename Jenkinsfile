pipeline {
  agent { label 'dagger' }

  environment {
    DAGGER_VERSION = "0.12.2"
    PATH = "/tmp/dagger/bin:$PATH"
    DOCKER_TOKEN = credentials('DOCKER_TOKEN')
    DAGGER_CLOUD_TOKEN =  credentials('DAGGER_CLOUD_TOKEN')
  }

  stages {
    stage("install dagger") {
      steps {
        sh '''
        curl -L https://dl.dagger.io/dagger/install.sh | BIN_DIR=/tmp/dagger/bin DAGGER_VERSION=$DAGGER_VERSION sh
        '''
      }
    }
    stage("lint") {
      steps {
        sh 'dagger call lint --dir . stdout'
      }
    }
    stage("test") {
      steps {
        sh 'dagger call test --dir . stdout'
      }
    }
    stage("publish") {
      steps {
        sh 'dagger call publish --dir . --token env:DOCKER_TOKEN'
      }
    }
  }
}