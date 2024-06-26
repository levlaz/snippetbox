pipeline {
  agent { label 'dagger' }

  environment {
    DOCKER_TOKEN = credentials('DOCKER_TOKEN')
    DAGGER_CLOUD_TOKEN =  credentials('DAGGER_CLOUD_TOKEN')
  }

  stages {
    stage("checkout") {
      steps {
        checkout ([
          changelog: false, 
          poll: false, 
          scm: scmGit(
            branches: [[name: '**']], 
            browser: github('https://github.com/levlaz/snippetbox'), 
            extensions: [
              cloneOption(
                honorRefspec: true, 
                noTags: true, 
                reference: '', 
                shallow: false
                ), 
                lfs(), 
                localBranch('main')
            ], 
            userRemoteConfigs: [
              [
                url: 'https://github.com/levlaz/snippetbox'
                ]
            ])])
      }
    }
    stage("install dagger") {
      steps {
        sh '''
        curl -L https://dl.dagger.io/dagger/install.sh | BIN_DIR=bin sh
        '''
      }
    }
    stage("lint") {
      steps {
        sh 'bin/dagger call lint --dir . stdout'
      }
    }
    stage("test") {
      steps {
        sh 'bin/dagger call test --dir . stdout'
      }
    }
    stage("publish") {
      steps {
        sh 'bin/dagger call publish --dir . --token env:DOCKER_TOKEN'
      }
    }
  }
}
