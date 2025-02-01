import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:firebase_messaging/firebase_messaging.dart';
import 'package:flutter_local_notifications/flutter_local_notifications.dart';
import 'package:uuid/uuid.dart';

void main() {
  runApp(MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Flutter Demo',
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.deepPurple),
        useMaterial3: true,
      ),
      home: DeviceTokenScreen(),
    );
  }
}

class DeviceTokenScreen extends StatefulWidget {
  const DeviceTokenScreen({super.key});

  @override
  DeviceTokenScreenState createState() => DeviceTokenScreenState();
}

class DeviceTokenScreenState extends State<DeviceTokenScreen> {
  String? _deviceToken;
  final FirebaseMessaging _firebaseMessaging = FirebaseMessaging.instance;
  final FlutterLocalNotificationsPlugin _flutterLocalNotificationsPlugin =
      FlutterLocalNotificationsPlugin();

  @override
  void initState() {
    super.initState();
    _loadDeviceToken();
    _setupFirebaseMessaging();
  }

  Future<void> _loadDeviceToken() async {
    SharedPreferences prefs = await SharedPreferences.getInstance();
    setState(() {
      _deviceToken = prefs.getString('deviceToken');
    });
  }

  Future<void> _generateDeviceToken() async {
    if (_deviceToken == null) {
      String newToken = const Uuid().v4();
      SharedPreferences prefs = await SharedPreferences.getInstance();
      await prefs.setString('deviceToken', newToken);
      setState(() {
        _deviceToken = newToken;
      });
    }
  }

  void _setupFirebaseMessaging() {
    _firebaseMessaging.requestPermission();
    FirebaseMessaging.onMessage.listen((RemoteMessage message) {
      RemoteNotification? notification = message.notification;
      AndroidNotification? android = message.notification?.android;
      if (notification != null && android != null) {
        _flutterLocalNotificationsPlugin.show(
          notification.hashCode,
          notification.title,
          notification.body,
          NotificationDetails(
            android: AndroidNotificationDetails(
              'channel_id',
              'channel_name',
              importance: Importance.max,
              priority: Priority.high,
            ),
          ),
        );
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('Device Token Information'),
      ),
      body: Center(
        child: _deviceToken != null
            ? Text('Your device token is: $_deviceToken')
            : ElevatedButton(
                onPressed: _generateDeviceToken,
                child: Text('Generate Device Token'),
              ),
      ),
    );
  }
}
