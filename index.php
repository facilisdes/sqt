<?php
const DB_HOST = 'mysql:3306';      /* Хост, к которому мы подключаемся */
const DB_USER = 'sqt_admin_1234';           /* Имя пользователя */
const DB_PASSWORD = 'P@ssw0rd-12POss@*';   /* Используемый пароль */
const DB_DATABASE = 'sqt';      /* База данных для запросов по умолчанию */
const DB_EVENTS_TABLE = 'events';      /* База данных для запросов по умолчанию */
const DB_USERS_TABLE = 'clients';      /* База данных для запросов по умолчанию */

const MD5_SALT_TEXT = "s0mE spicy_s@lt h3re";

session_start();
session_regenerate_id();

function closeForNonAdmin() {
    if(empty($_SESSION['user_client']) || $_SESSION['user_client'] != 'admin')
        die();
}

if($_POST['ajax'] =='Y') {

}
else {
    header('Content-Type: text/html; charset=UTF-8');
    if($_REQUEST['mode'] == 'logout') {
        unset($_SESSION['user_id']);
        unset($_SESSION['user_client']);
        header("Location: /sqt.php");
        die();
    }
    ?>
    <html>
    <head>
        <title>Dashboard</title>
    </head>
    <body>
    <table>
        <tr>
            <?php if(empty($_SESSION['user_client'])):?>
                <td><a href="?mode=auth" <?php if($_REQUEST['mode'] == 'auth') echo 'style="font-weight: 900;"' ?> >Авторизация</a></td>
                <td><a href="?mode=client_add" <?php if($_REQUEST['mode'] == 'client_add') echo 'style="font-weight: 900;"' ?> >Добавление клиента</a></td>
            <?php elseif($_SESSION['user_client'] == 'admin'):?>
                <td><a href="?mode=list_select" <?php if($_REQUEST['mode'] == 'list_select') echo 'style="font-weight: 900;"' ?> >Общая статистика</a></td>
                <td><a href="?mode=all" <?php if($_REQUEST['mode'] == 'all') echo 'style="font-weight: 900;"' ?> >Все запросы</a></td>
                <td><a href="?mode=clients" <?php if($_REQUEST['mode'] == 'clients') echo 'style="font-weight: 900;"' ?> >Клиенты</a></td>
                <td><a href="?mode=client_add" <?php if($_REQUEST['mode'] == 'client_add') echo 'style="font-weight: 900;"' ?> >Добавление клиента</a></td>
                <td><a href="?mode=logout" <?php if($_REQUEST['mode'] == 'logout') echo 'style="font-weight: 900;"' ?> >Выход</a></td>
            <?php else:?>
                <td><a href="?mode=list" <?php if($_REQUEST['mode'] == 'list') echo 'style="font-weight: 900;"' ?> >Статистика</a></td>
                <td><a href="?mode=logout" <?php if($_REQUEST['mode'] == 'logout') echo 'style="font-weight: 900;"' ?> >Выход</a></td>
            <?php endif;?>


            <td><?php echo $_SESSION['user_client']; ?></td>
        </tr>
    </table>
    <?php switch($_REQUEST['mode']) {?>
<?php case 'list':?>
        <?php
        $arRows = [];
        $arRowsStatuses = [];
        $successCount = 0;
        $failureCount = 0;

        if(!empty($_SESSION['user_client'])) {

            $link = mysqli_connect(DB_HOST, DB_USER, DB_PASSWORD, DB_DATABASE);

            if (!$link) {
                die("Невозможно подключиться к базе данных: " . mysqli_connect_error());
            }

            $query = 'SELECT AVG(TimeElapsed) AS te, AVG(TimeElapsedTotal) AS tet FROM ' . DB_EVENTS_TABLE . " WHERE Client = '".$_SESSION['user_client']."' AND IsExecuted = 1 GROUP BY Client;";

            if ($result = mysqli_query($link, $query)) {
                while ($row = mysqli_fetch_assoc($result)) {
                    $arRows[] = [
                        'Client' => $_SESSION['user_client'],
                        'te' => round(floatval($row['te']), 3),
                        'tet' => round(floatval($row['tet']), 3),
                    ];
                }
                mysqli_free_result($result);
            }

            $query = 'SELECT COUNT(IsExecuted) AS successCount FROM ' . DB_EVENTS_TABLE . " WHERE Client = '".$_SESSION['user_client']."' AND IsExecuted = 1 GROUP BY Client;";

            if ($result = mysqli_query($link, $query)) {
                while ($row = mysqli_fetch_assoc($result)) {
                    $successCount = $row['successCount'];
                }
                mysqli_free_result($result);
            }

            $query = 'SELECT COUNT(IsExecuted) AS failureCount FROM ' . DB_EVENTS_TABLE . " WHERE Client = '".$_SESSION['user_client']."' AND IsExecuted = 0 GROUP BY Client;";

            if ($result = mysqli_query($link, $query)) {
                while ($row = mysqli_fetch_assoc($result)) {
                    $failureCount = $row['failureCount'];
                }
                mysqli_free_result($result);
            }

            $query = "SELECT COUNT(StatusText) as ct, StatusText as st FROM ".DB_EVENTS_TABLE." WHERE Client='".$_SESSION['user_client']."' GROUP BY StatusText;";

            if ($result = mysqli_query($link, $query)) {
                while ($row = mysqli_fetch_assoc($result)) {
                    $arRowsStatuses[] = [
                        'status' => $row['st'],
                        'count' => $row['ct'],
                    ];
                }
                mysqli_free_result($result);
            }


            mysqli_close($link);
        }
        ?>
        <table border="1" cellpadding="1" cellspacing="1" style="table-layout: fixed;">
            <tr>
                <td colspan="2">Данные по успешным запросам</td>
            </tr>
            <tr>
                <td>
                    Среднее время запросов без учёта очереди, мс
                </td>
                <td>
                    Среднее время запросов с очередью, мс
                </td>
            </tr>
            <?php foreach($arRows as $arRow):?>
                <tr>
                    <td>
                        <?php echo $arRow['te'] ?>
                    </td>
                    <td>
                        <?php echo $arRow['tet'] ?>
                    </td>
                </tr>
            <?endforeach;?>
        </table>
        <br>
        <table border="1" cellpadding="1" cellspacing="1" style="table-layout: fixed;">
            <tr>
                <td colspan="2">Данные по результативности запросов</td>
            </tr>
            <tr>
                <td>
                    Успешно
                </td>
                <td>
                    Неудачно
                </td>
            </tr>
            <tr>
                <td>
                    <?php echo $successCount ?>
                </td>
                <td>
                    <?php echo $failureCount ?>
                </td>
            </tr>
        </table>
        <br>
        <table border="1" cellpadding="1" cellspacing="1" style="table-layout: fixed;">
            <tr>
                <td colspan="2">Данные по статусам</td>
            </tr>
            <tr>
                <td>
                    Статус
                </td>
                <td>
                    Количество запросов
                </td>
            </tr>
            <?php foreach($arRowsStatuses as $arRow):?>
                <tr>
                    <td>
                        <?php echo $arRow['status'] ?>
                    </td>
                    <td>
                        <?php echo $arRow['count'] ?>
                    </td>
                </tr>
            <?endforeach;?>
        </table>
        <br>
<?php break;
case 'list_select':?>
        <?php closeForNonAdmin();?>
        <?php
        $arRows = [];
        $arRowsStatuses = [];
        $arClients = [];
        $successCount = 0;
        $failureCount = 0;

        $client = urldecode($_POST['client']);


        $link = mysqli_connect(DB_HOST, DB_USER, DB_PASSWORD, DB_DATABASE);

        if (!$link) {
            die("Невозможно подключиться к базе данных: " . mysqli_connect_error());
        }

        $query = 'SELECT * FROM ' . DB_USERS_TABLE . ";";
        if ($result = mysqli_query($link, $query)) {
            while( $row = mysqli_fetch_assoc($result) ){
                if(empty($client)) $client = $row['address'];
                $arClients[] = [
                    'ID' => $row['id'],
                    'address' => $row['address'],
                    'note' => $row['note']
                ];
            }
            mysqli_free_result($result);
        }?>
        <form method="post">
            <input type="hidden" name="mode" value="list_select">
            <label for="client">Клиент:</label>
            <select id="client" name="client" >
                <?php foreach($arClients as $arClient):?>
                    <option value="<?php echo $arClient['address'];?>" <?php if($client == $arClient['address']) echo 'selected';?> >
                        <?php echo $arClient['address'];?>
                    </option>
                <?php endforeach;?>
            </select>
            <br>
            <input type="submit">
        </form>

        <?php
        if(!empty($client)) {

            $query = 'SELECT AVG(TimeElapsed) AS te, AVG(TimeElapsedTotal) AS tet FROM ' . DB_EVENTS_TABLE . " WHERE Client = '".$client."' AND IsExecuted = 1 GROUP BY Client;";

            if ($result = mysqli_query($link, $query)) {
                while ($row = mysqli_fetch_assoc($result)) {
                    $arRows[] = [
                        'Client' => $client,
                        'te' => round(floatval($row['te']), 3),
                        'tet' => round(floatval($row['tet']), 3),
                    ];
                }
                mysqli_free_result($result);
            }

            $query = 'SELECT COUNT(IsExecuted) AS successCount FROM ' . DB_EVENTS_TABLE . " WHERE Client = '".$client."' AND IsExecuted = 1 GROUP BY Client;";

            if ($result = mysqli_query($link, $query)) {
                while ($row = mysqli_fetch_assoc($result)) {
                    $successCount = $row['successCount'];
                }
                mysqli_free_result($result);
            }

            $query = 'SELECT COUNT(IsExecuted) AS failureCount FROM ' . DB_EVENTS_TABLE . " WHERE Client = '".$client."' AND IsExecuted = 0 GROUP BY Client;";

            if ($result = mysqli_query($link, $query)) {
                while ($row = mysqli_fetch_assoc($result)) {
                    $failureCount = $row['failureCount'];
                }
                mysqli_free_result($result);
            }

            $query = "SELECT COUNT(StatusText) as ct, StatusText as st FROM ".DB_EVENTS_TABLE." WHERE Client='".$client."' GROUP BY StatusText;";

            if ($result = mysqli_query($link, $query)) {
                while ($row = mysqli_fetch_assoc($result)) {
                    $arRowsStatuses[] = [
                        'status' => $row['st'],
                        'count' => $row['ct'],
                    ];
                }
                mysqli_free_result($result);
            }


            mysqli_close($link);
        }
        ?>
        <table border="1" cellpadding="1" cellspacing="1" style="table-layout: fixed;">
            <tr>
                <td colspan="2">Данные по успешным запросам</td>
            </tr>
            <tr>
                <td>
                    Среднее время запросов без учёта очереди, мс
                </td>
                <td>
                    Среднее время запросов с очередью, мс
                </td>
            </tr>
            <?php foreach($arRows as $arRow):?>
                <tr>
                    <td>
                        <?php echo $arRow['te'] ?>
                    </td>
                    <td>
                        <?php echo $arRow['tet'] ?>
                    </td>
                </tr>
            <?endforeach;?>
        </table>
        <br>
        <table border="1" cellpadding="1" cellspacing="1" style="table-layout: fixed;">
            <tr>
                <td colspan="2">Данные по результативности запросов</td>
            </tr>
            <tr>
                <td>
                    Успешно
                </td>
                <td>
                    Неудачно
                </td>
            </tr>
            <tr>
                <td>
                    <?php echo $successCount ?>
                </td>
                <td>
                    <?php echo $failureCount ?>
                </td>
            </tr>
        </table>
        <br>
        <table border="1" cellpadding="1" cellspacing="1" style="table-layout: fixed;">
            <tr>
                <td colspan="2">Данные по статусам</td>
            </tr>
            <tr>
                <td>
                    Статус
                </td>
                <td>
                    Количество запросов
                </td>
            </tr>
            <?php foreach($arRowsStatuses as $arRow):?>
                <tr>
                    <td>
                        <?php echo $arRow['status'] ?>
                    </td>
                    <td>
                        <?php echo $arRow['count'] ?>
                    </td>
                </tr>
            <?endforeach;?>
        </table>
        <br>

<?php break;
case 'auth':?>
        <?php
        $showForm = false;
        $note = false;
        if(empty($_POST['address']) || empty($_POST['passwd']))
            $showForm = true;

        if(!$showForm) {

            $address = urldecode($_POST['address']);
            $password = urldecode($_POST['passwd']);

            $passwordForDb = md5(MD5_SALT_TEXT . $password . $address);

            $link = mysqli_connect(DB_HOST, DB_USER, DB_PASSWORD, DB_DATABASE);

            if (!$link) {
                die("Невозможно подключиться к базе данных: " . mysqli_connect_error());
            }

            $query = "SELECT id FROM " . DB_USERS_TABLE . " WHERE address = ? AND password = ? "; //

            if($stmt = mysqli_prepare($link, $query)) {
                mysqli_stmt_bind_param($stmt, "ss", $address, $passwordForDb);

                mysqli_stmt_execute($stmt);

                $res = mysqli_stmt_get_result($stmt);

                $id = -1;
                if($res && $row = $res->fetch_assoc()) {
                    $id = $row['id'];
                }
                if(empty($id)) {
                    $showForm = true;
                    $note = "Ошибка авторизации! Неверный логин или пароль.";
                }
                else {
                    $_SESSION['user_id'] = $id;
                    $_SESSION['user_client'] = $address;
                }
            }

            $note = "Успешная авторизация!";
            mysqli_close($link);
        }
        ?>
        <?php if($showForm):?>
        <form method="post">
            <input type="hidden" name="mode" value="auth">
            <label for="address">Адрес сервера:</label>
            <input type="text" id="address" name="address" placeholder="127.0.0.1" value="<?=$_POST['address']?>">
            <br>
            <label for="address">Пароль:</label>
            <input type="password" id="passwd" name="passwd" value="<?=$_POST['passwd']?>">
            <br>
            <input type="submit">
        </form>
    <?endif;?>
        <p><?php echo $note;?></p>
    <?php if(!empty($_SESSION['user_id'])):?>
        <script>window.location.replace("/sqt.php");</script>
    <?endif;?>

<?php break;
case 'all':?>
        <?php closeForNonAdmin();?>
        <?php
        $client = urldecode($_REQUEST['filter_client']);
        if(empty($client)) $client = 'all';

        $link = mysqli_connect(DB_HOST, DB_USER, DB_PASSWORD, DB_DATABASE);

        if (!$link) {
            die("Невозможно подключиться к базе данных: " . mysqli_connect_error());
        }
        $arClients = [];
        $query = 'SELECT * FROM ' . DB_USERS_TABLE . ";";
        if ($result = mysqli_query($link, $query)) {
            while( $row = mysqli_fetch_assoc($result) ){
                $arClients[] = [
                    'ID' => $row['id'],
                    'address' => $row['address'],
                    'note' => $row['note']
                ];
            }
            mysqli_free_result($result);
        }


        ?>
        <form method="get">
            <input type="hidden" name="mode" value="all">
            <label for="filter_client">Фильтрация по клиенту:</label>
            <select id="filter_client" name="filter_client" >
                <option value="<?php echo $arClient['address'];?>" <?php if($client == 'all') echo 'selected';?> >
                    Все клиенты
                </option>
                <?php foreach($arClients as $arClient):?>
                    <option value="<?php echo $arClient['address'];?>" <?php if($client == $arClient['address']) echo 'selected';?> >
                        <?php echo $arClient['address'];?>
                    </option>
                <?php endforeach;?>
            </select>
            <input type="submit">
        </form>
        <?php
        $arRows = [];

        $query = 'SELECT * FROM ' . DB_EVENTS_TABLE . ' ';
        $where = false;
        if(!empty($client) && $client != 'all') {
            $where = "%" . $client . "%";
            $query .= "WHERE Client LIKE ? "; //
        }
        $query.='ORDER BY id ASC;';

        if($stmt = mysqli_prepare($link, $query)) {
            if($where)
                mysqli_stmt_bind_param($stmt, "s", $where);

            mysqli_stmt_execute($stmt);

            $res = mysqli_stmt_get_result($stmt);

            $arCommands = [
                0 => 'Обычный запрос',
                1 => 'Healthcheck'
            ];
            while($res && $row = $res->fetch_assoc()) {
                $arRows[] = [
                    'ID' => $row['id'],
                    'Client' => $row['Client'],
                    'IsExecuted' => $row['IsExecuted'] == 1 ? "Да" : "Нет",
                    'Status' => $row['Status'],
                    'StatusText' => $row['StatusText'],
                    'RequestedKey' => $row['RequestedKey'],
                    'Data' => $row['Data'],
                    'LocalData' => $row['LocalData'],
                    'ValueIsValidated' => $row['ValueIsValidated'] == 1 ? "Да" : "Нет",
                    'TimeStart' => $row['TimeStart'] ? date("m.d H:i:s", $row['TimeStart']) : "",
                    'TimeEnd' => $row['TimeEnd'] ? date("m.d H:i:s", $row['TimeEnd']) : "",
                    'TimeElapsed' => $row['TimeElapsed'],
                    'TimeElapsedTotal' => $row['TimeElapsedTotal'],
                    'TimeQueuedMin' => $row['TimeQueuedMin'],
                    'QueueSize' => $row['QueueSize'],
                    'Command' => empty($arCommands[$row['Command']]) ? '-' : $arCommands[$row['Command']],
                ];
            }
        }
        mysqli_close($link);
        ?>
        <div id="wrap">
            <table border="1" cellpadding="1" cellspacing="1" style="table-layout: fixed;">
                <thead >
                <tr style="display:block;">
                    <td style="width:4%">ID записи</td>
                    <td style="width:9%">Клиент</td>
                    <td style="width:5%">Выполнен ли</td>
                    <td style="width:4%">ID статуса</td>
                    <td style="width:8%">Статус</td>
                    <td style="width:6%">Запрос</td>
                    <td style="width:10%">Ответ</td>
                    <td style="width:10%">Ответ (локальное хранилище)</td>
                    <td style="width:4%">Ответы совпали</td>
                    <td style="width:6%">Начало исполнения</td>
                    <td style="width:6%">Конец исполнения</td>
                    <td style="width:6%">Затрачено на запрос, мс</td>
                    <td style="width:6%">Затрачено на запрос + очередь, мс</td>
                    <td style="width:6%">Мин. время нахождения в очереди, мс</td>
                    <td style="width:5%">Размер очереди (вместе с данным запросом)</td>
                    <td style="width:5%">Команда</td>
                </tr>
                </thead>
                <tbody style="display:block;overflow:auto;height:75vh;width:100%;">
                <?php foreach($arRows as $arRow):?>
                    <tr>
                        <td style="width:4%"><?php echo $arRow['ID']; ?></td>
                        <td style="width:9%"><?php echo $arRow['Client']; ?></td>
                        <td style="width:5%"><?php echo $arRow['IsExecuted']; ?></td>
                        <td style="width:4%"><?php echo $arRow['Status']; ?></td>
                        <td style="width:8%"><?php echo $arRow['StatusText']; ?></td>
                        <td style="width:6%"><?php echo $arRow['RequestedKey']; ?></td>
                        <td style="width:10%"><?php echo $arRow['Data']; ?></td>
                        <td style="width:10%"><?php echo $arRow['LocalData']; ?></td>
                        <td style="width:4%"><?php echo $arRow['ValueIsValidated']; ?></td>
                        <td style="width:6%"><?php echo $arRow['TimeStart']; ?></td>
                        <td style="width:6%"><?php echo $arRow['TimeEnd']; ?></td>
                        <td style="width:6%"><?php echo $arRow['TimeElapsed']; ?></td>
                        <td style="width:6%"><?php echo $arRow['TimeElapsedTotal']; ?></td>
                        <td style="width:6%"><?php echo $arRow['TimeQueuedMin']; ?></td>
                        <td style="width:5%"><?php echo $arRow['QueueSize']; ?></td>
                        <td style="width:5%"><?php echo $arRow['Command']; ?></td>
                    </tr>
                <?php endforeach; ?>
                </tbody>
            </table>
        </div>

<?php break;
case 'clients':?>
        <?php closeForNonAdmin();?>
        <?php
        $arRows = [];

        $link = mysqli_connect(DB_HOST, DB_USER, DB_PASSWORD, DB_DATABASE);

        if (!$link) {
            die("Невозможно подключиться к базе данных: " . mysqli_connect_error());
        }

        $query = 'SELECT * FROM ' . DB_USERS_TABLE . ";";

        if ($result = mysqli_query($link, $query)) {
            while( $row = mysqli_fetch_assoc($result) ){
                $arRows[] = [
                    'ID' => $row['id'],
                    'address' => $row['address'],
                    'note' => $row['note']
                ];
            }
            mysqli_free_result($result);
        }
        mysqli_close($link);
        ?>
        <table border="1" cellpadding="1" cellspacing="1" style="table-layout: fixed;">
            <tr>
                <td>ID</td>
                <td>Адрес</td>
                <td>Примечание</td>
            </tr>
            <?php foreach($arRows as $arRow):?>
                <tr>
                    <td><?php echo $arRow['ID']?></td>
                    <td><?php echo $arRow['address']?></td>
                    <td><?php echo $arRow['note']?></td>
                </tr>
            <?php endforeach;?>
        </table>

<?php break;
case 'client_add':?>
    <?php if(empty($_POST['address']) || empty($_POST['note'])):?>
        <form method="post">
            <input type="hidden" name="mode" value="client_add">
            <label for="address">Адрес сервера:</label>
            <input type="text" id="address" name="address" placeholder="127.0.0.1" value="<?=$_POST['address']?>">
            <br>
            <label for="address">Описание:</label>
            <input type="text" id="note" name="note" value="<?=$_POST['note']?>">
            <br>
            <label for="address">Пароль для входа:</label>
            <input type="password" id="passwd" name="passwd" value="<?=$_POST['passwd']?>">
            <br>
            <input type="submit">
        </form>
    <?php else:
    $address = urldecode($_POST['address']);
    $note = urldecode($_POST['note']);
    $password = utf8_decode(urldecode($_POST['passwd']));

    $passwordForDb = md5(MD5_SALT_TEXT . $password . $address);

    $link = mysqli_connect(DB_HOST, DB_USER, DB_PASSWORD, DB_DATABASE);

    if (!$link) {
        die("Невозможно подключиться к базе данных. Код ошибки: %s\n" . mysqli_connect_error());
    }

    $query = "INSERT INTO " . DB_USERS_TABLE . " (address, note, password) VALUES (?, ?, ?);";

    $stmt = mysqli_prepare($link, $query);
    mysqli_stmt_bind_param($stmt, "sss", $address, $note, $password);



    if ($res = mysqli_stmt_execute($stmt)):?>
        <p>Данные успешно сохранены!</p>
    <?php else:?>
    <p>Ошибка при сохранении: <?php echo mysqli_error($link) ?>
        <?php endif;
        mysqli_close($link);
        ?>
        <?php endif;?>
        <?php }?>
    </body>
    </html>

    <?php
}
