use 5.014;
use strict;
use warnings;


sub main {
	my $name = "deploy";
	system("go install $name");
	die if $?;
	system("bin\\$name");
}

#
main