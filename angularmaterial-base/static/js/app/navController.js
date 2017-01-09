(function () {

    angular.module('baseApp')
        .controller('navController', ['$mdSidenav', navController]);
    // $mdSidenav is dependency injected into the controller when controller is built. This is done by name
    // The reason why it's included as a string above is so that you can minify it. Check out one of my controller classes for reference
    function navController($mdSidenav) {
        var self = this;

        self.message = 'Sup yeow';
        self.buttonText = 'Click me!';
        self.auxButtons = ['lol', 'button kek'];

        self.sayMessage = function () {
            self.message = 'Seriously though, get bower... Look at the dependency list for mine below. All managed automagically';

            self.otherStyle = {
                color: 'red',
                'font-size': '14px'
            };
            
            self.otherMessage = 'dependencies for mine: "bower install" on server and i have all of them up to date in the proper place\n';
        };

        self.toggle = function () {
            $mdSidenav('left').toggle();
        };
        
        /* called when an aux button is pressed */
        self.selectAux = function(name) {
			alert(name);
		};
    }
})();
