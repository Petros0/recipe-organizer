import 'dart:developer';

import 'package:appwrite/appwrite.dart';
import 'package:flutter/material.dart';
import 'package:forui/forui.dart';
import 'package:recipe_organizer/core/di.dart';
import 'package:signals/signals_flutter.dart';

class CounterPage extends StatefulWidget {
  const CounterPage({super.key});

  @override
  State<CounterPage> createState() => _CounterPageState();
}

class _CounterPageState extends State<CounterPage> {
  final FlutterSignal<int> _counter = signal<int>(0);

  void _increment() => _counter.value++;
  void _decrement() => _counter.value--;

  Future<void> _sendPing() async {
    try {
      await getIt<Client>().ping();
    } on AppwriteException catch (e) {
      log('Ping failed: $e');
    }
  }

  @override
  Widget build(BuildContext context) {
    return FScaffold(
      header: FHeader(
        title: const Text('Settings'),
        suffixes: [FHeaderAction(icon: const Icon(FIcons.ellipsis), onPress: () {})],
      ),
      child: Column(
        children: [
          Text('Counter: ${_counter.watch(context)}'),
          FButton(onPress: _increment, child: const Icon(Icons.add)),
          FButton(onPress: _decrement, child: const Icon(Icons.remove)),
          const SizedBox(height: 12),
          FButton(onPress: _sendPing, child: const Text('Send a ping')),
        ],
      ),
    );
  }
}
