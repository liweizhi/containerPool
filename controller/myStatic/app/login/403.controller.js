(function(){
	'use strict';

	angular
	    .module('container-pool.login')
	    .controller('AccessDeniedController', AccessDeniedController);

	AccessDeniedController.$inject = ['$stateParams'];
	function AccessDeniedController($stateParams) {
            var vm = this;
	}
})();
