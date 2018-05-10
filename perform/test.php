#!php
<?php
$a = 1;

$start = microtime(true);
for ($i = 1; $i < 100000000; $i++) {
    $a += $i;
    $a % $i;
}
$elapsed = microtime(true) - $start;
echo "PHP took {$elapsed} seconds\n";
