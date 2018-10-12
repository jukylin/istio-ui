<?php
define('CURRENT_TIMESTAMP', $_SERVER['REQUEST_TIME']);
define('CURRENT_DATETIME', date('Y-m-d H:i:s'));
define('COOKIE_PATH', '/');
define('COOKIE_DOMAIN', '.etcchebao.com');
define('COOKIE_PREFIX', 'MB_');
define('COOKIE_EXPIRE', 7);

define('DS', DIRECTORY_SEPARATOR);
define('PS', PATH_SEPARATOR);
define('WEB_DIR', dirname(__FILE__) );

if (isset($_SERVER['MB_APPLICATION']) && $_SERVER['MB_APPLICATION'] == 'production') {
	defined('YII_DEBUG') or define('YII_DEBUG', false);
	defined('YII_ENV') or define( 'YII_ENV', 'pro');
	defined('YII_TRACE_LEVEL') or define('YII_TRACE_LEVEL', 1);
	define('VDIR', '/home/wwwroot/');
	$config = require(__DIR__ . '/../config/production.php');
} else {
	ini_set('display_errors', 'on');
	error_reporting(E_ALL);
	defined('YII_DEBUG') or define('YII_DEBUG', true);
	defined('YII_ENV') or define( 'YII_ENV', 'dev');
	defined('YII_TRACE_LEVEL') or define('YII_TRACE_LEVEL', 3);
	define('VDIR', dirname(dirname(__DIR__)));
	$config = require(__DIR__ . '/../config/development.php');
}

require(VDIR . '/framework/vendor/autoload.php');
require(VDIR . '/framework/vendor/yiisoft/yii2/Yii.php');
# 目录映射
Yii::setAlias('@BaseComponents', VDIR . '/framework/components/');
(new yii\web\Application($config))->run();
