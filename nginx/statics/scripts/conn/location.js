let locationInterval = null;

// 3秒ごとに位置情報をポストする
locationInterval = setInterval(async () => {
    // jwt を取得
    const jwtToken = await GetJwt();

    // 現在位置を取得
    const location = GetLocation();

    // post で送信
    const req = await fetch("/location/update",{
        method: "POST",
        headers: {
            "Content-type": "application/json",
            "Authorization": "Bearer " + jwtToken
        },
        body: JSON.stringify({
            "lat" : location.latitude,
            "lng" : location.longitude,
        })
    });

    const result = await req.json();
    result["near"].forEach(nUser => {
        // ピンが存在するか判定
        if (ExistPin(nUser["userid"])) {
            // 存在するとき動かす
            MovePin(nUser["userid"],nUser["latitude"],nUser["longitude"]);
            return;
        }

        // 近くにいるフレンド表示
        NewPin(nUser["userid"],GetIcon(nUser["userid"]),nUser["latitude"],nUser["longitude"]);
    });

    // 離れたフレンドのピンを消す
    result["removed"].forEach(dUserID => {
        // ピンが存在するか判定
        if (ExistPin(dUserID)) {
            // 存在するとき消す
            RemovePin(dUserID);

            return;
        }
    })
}, 2000);