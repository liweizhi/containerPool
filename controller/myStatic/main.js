/**
 * Created by james on 17/5/9.
 */

(function(){
    var app = angular.module('demo', ['ui.router']);
    app.config(['$stateProvider', '$urlRouterProvider', function($stateProvider, $urlRouterProvider) {
        $urlRouterProvider.otherwise(function ($injector) {
            var $state = $injector.get('$state');
            console.log("other");
            $state.go('shows.detail');
        });

        $stateProvider
            .state('shows', {
                url:'/hello',
                templateUrl: 'shows.html',
                controller: 'ShowsController'
            })
            .state('shows.detail', {
                url: '^/detail/:id',
                templateUrl: './show-detail.html',
                controller: 'ShowsDetailController'
            });
    }]);
    app.factory('ShowsService',function(){
        var shows = [{
            id: 1,
            name: 'Walking Dead',
            description: 'The Walking Dead is an American post-apocalyptic horror drama television series developed by Frank Darabont. It is based on the comic book series of the same name by Robert Kirkman, Tony Moore, and Charlie Adlard. It stars Andrew Lincoln as sheriff\'s deputy Rick Grimes, who awakens from a coma to find a post-apocalyptic world dominated by flesh-eating zombies.'
        },
            {
                id: 2,
                name: 'Breaking Bad',
                description: 'Breaking Bad is an American crime drama television series created and produced by Vince Gilligan. The show originally aired on the AMC network for five seasons, from January 20, 2008 to September 29, 2013. The main character is Walter White (Bryan Cranston), a struggling high school chemistry teacher who is diagnosed with inoperable lung cancer at the beginning of the series.'
            },
            {
                id: 3,
                name: '7D',
                description: 'The 7D is an American animated television series produced by Disney Television Animation, and broadcast on Disney XD starting in July 7, 2014. It is a re-imagining of the titular characters from the 1937 film Snow White and the Seven Dwarfs by Walt Disney Productions'
            }];


        return {
            list: function(){
                return shows;
            },
            find: function(id){
                return _.find(shows, function(show){return show.id == id});
            }
        }
    });
    app.controller('ShowsController', ['$scope','ShowsService', function($scope, ShowsService) {
        $scope.shows = ShowsService.list();
    }]);

    app.controller('ShowsDetailController', ['$scope','$stateParams', 'ShowsService', function($scope, $stateParams, ShowsService) {
        $scope.selectedShow = ShowsService.find($stateParams.id);
    }]);


})();
