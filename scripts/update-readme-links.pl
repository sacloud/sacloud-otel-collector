#!/usr/bin/env perl
# Rewrite component documentation links in README.md so that they point to
# the exact versions listed in builder-config.yaml.
#
# Usage:
#   scripts/update-readme-links.pl          # update README.md in place
#   scripts/update-readme-links.pl --check  # exit 1 if README.md is out of date
use strict;
use warnings;

my $check  = @ARGV && $ARGV[0] eq '--check';
my $config = 'builder-config.yaml';
my $readme = 'README.md';

# Vanity import paths that do not map directly to github.com/<owner>/<repo>.
my %vanity = (
    'go.opentelemetry.io/collector' => 'open-telemetry/opentelemetry-collector',
);

# Collect component modules from builder-config.yaml.
# Providers are not listed in README, so they are skipped.
my @modules;
my $section = '';
open my $cfg, '<', $config or die "cannot open $config: $!\n";
while (<$cfg>) {
    if (/^(\w+):/) {
        $section = $1;
        next;
    }
    next unless $section =~ /^(?:receivers|processors|exporters|extensions|connectors)$/;
    next unless /^\s*-\s*gomod:\s*(\S+)\s+(\S+)/;
    my ($mod, $ver) = ($1, $2);
    next if $mod =~ m{^github\.com/sacloud/sacloud-otel-collector/};    # local modules

    my ($repo, $sub);
    for my $prefix (keys %vanity) {
        if ($mod =~ m{^\Q$prefix\E/(.+)$}) {
            ($repo, $sub) = ($vanity{$prefix}, $1);
        }
    }
    if (!$repo && $mod =~ m{^github\.com/([^/]+/[^/]+)/(.+)$}) {
        ($repo, $sub) = ($1, $2);
    }
    die "$config: cannot map module $mod to a GitHub URL\n" unless $repo;

    # Go modules in subdirectories are tagged as <subdir>/<version>.
    push @modules, {
        mod  => $mod,
        repo => $repo,
        sub  => $sub,
        url  => "https://github.com/$repo/tree/$sub/$ver/$sub",
    };
}
close $cfg;

open my $in, '<', $readme or die "cannot open $readme: $!\n";
my $orig = do { local $/; <$in> };
close $in;

my $content = $orig;
my @errors;
my %expected;
for my $m (@modules) {
    my $n = $content =~ s{https://github\.com/\Q$m->{repo}\E/tree/[^)\s]+?/\Q$m->{sub}\E(?=\))}{$m->{url}}g;
    push @errors, "$readme: no documentation link for $m->{mod}" unless $n;
    $expected{ $m->{url} } = 1;
}

# Every GitHub tree link must correspond to a module in builder-config.yaml.
while ($content =~ m{(https://github\.com/[^/\s]+/[^/\s]+/tree/[^)\s]+)}g) {
    push @errors, "$readme: $1 does not match any module in $config" unless $expected{$1};
}

if (@errors) {
    print STDERR "$_\n" for @errors;
    exit 1;
}

if ($content eq $orig) {
    exit 0;
}
if ($check) {
    print STDERR "$readme: component links are out of date. Run `make readme-links`.\n";
    exit 1;
}
open my $out, '>', $readme or die "cannot write $readme: $!\n";
print $out $content;
close $out;
print "updated $readme\n";
