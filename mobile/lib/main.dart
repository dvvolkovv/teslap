import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import 'app.dart';
import 'core/di/injection.dart';

/// Application entry point.
///
/// Initialises dependency injection, locks orientation to portrait,
/// and launches the TeslaPay app.
void main() async {
  WidgetsFlutterBinding.ensureInitialized();

  // Lock to portrait orientation for phone-first experience.
  await SystemChrome.setPreferredOrientations([
    DeviceOrientation.portraitUp,
    DeviceOrientation.portraitDown,
  ]);

  // System UI overlay styling.
  SystemChrome.setSystemUIOverlayStyle(
    const SystemUiOverlayStyle(
      statusBarColor: Colors.transparent,
      statusBarIconBrightness: Brightness.dark,
      systemNavigationBarColor: Colors.white,
      systemNavigationBarIconBrightness: Brightness.dark,
    ),
  );

  // Dependency injection.
  await configureDependencies();

  runApp(const TeslaPayApp());
}
